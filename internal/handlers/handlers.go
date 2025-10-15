package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"

	"chat-server/internal/config"
	"chat-server/internal/message"
	"chat-server/internal/room"
	"chat-server/internal/user"
)

type Handlers struct {
	roomManager *room.Manager
	userManager *user.Manager
	config      *config.Config
	upgrader    websocket.Upgrader
}

func New(roomManager *room.Manager, userManager *user.Manager, config *config.Config) *Handlers {
	return &Handlers{
		roomManager: roomManager,
		userManager: userManager,
		config:      config,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}

func (h *Handlers) ServeHome(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/index.html")
}

func (h *Handlers) GetRooms(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	rooms := h.roomManager.GetRoomsInfo()

	response := map[string]interface{}{
		"success": true,
		"rooms":   rooms,
	}

	json.NewEncoder(w).Encode(response)
}

func (h *Handlers) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Ошибка обновления до WebSocket: %v", err)
		return
	}
	defer conn.Close()

	username := r.URL.Query().Get("username")
	roomName := r.URL.Query().Get("room")

	if err := h.validateConnection(username, roomName); err != nil {
		h.sendError(conn, err.Error())
		return
	}

	userID := uuid.New().String()
	user := user.NewUser(userID, username, conn)
	h.userManager.AddUser(user)
	defer h.userManager.RemoveUser(userID)

	roomID := h.generateRoomID(roomName)
	chatRoom := h.roomManager.GetOrCreateRoom(roomID, roomName)
	chatRoom.AddUser(user)
	defer chatRoom.RemoveUser(user)

	h.broadcastUserJoin(chatRoom, user)

	h.sendMessageHistory(conn, chatRoom)

	h.handleMessages(user, chatRoom)
}

func (h *Handlers) validateConnection(username, roomName string) error {
	if username == "" {
		return fmt.Errorf("имя пользователя обязательно")
	}

	if len(username) > h.config.MaxUsernameLength {
		return fmt.Errorf("имя пользователя слишком длинное (максимум %d символов)", h.config.MaxUsernameLength)
	}

	if matched, _ := regexp.MatchString(`[<>'"&]`, username); matched {
		return fmt.Errorf("имя пользователя содержит недопустимые символы")
	}

	if roomName == "" {
		return fmt.Errorf("название комнаты обязательно")
	}

	if len(roomName) > h.config.MaxRoomNameLength {
		return fmt.Errorf("название комнаты слишком длинное (максимум %d символов)", h.config.MaxRoomNameLength)
	}

	return nil
}

func (h *Handlers) generateRoomID(roomName string) string {
	return strings.ReplaceAll(strings.ToLower(roomName), " ", "-")
}

func (h *Handlers) handleMessages(u *user.User, room *room.Room) {
	conn := u.Conn
	conn.SetReadLimit(512)
	conn.SetReadDeadline(time.Now().Add(time.Duration(h.config.PongWait) * time.Second))
	conn.SetPongHandler(func(string) error {
		u.UpdatePing()
		conn.SetReadDeadline(time.Now().Add(time.Duration(h.config.PongWait) * time.Second))
		return nil
	})

	ticker := time.NewTicker(time.Duration(h.config.PingInterval) * time.Second)
	defer ticker.Stop()

	go func() {
		for range ticker.C {
			conn.SetWriteDeadline(time.Now().Add(time.Duration(h.config.WriteWait) * time.Second))
			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}()

	for {
		var msg message.MessageRequest
		err := conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket ошибка: %v", err)
			}
			break
		}

		h.processMessage(u, room, &msg)
	}
}

func (h *Handlers) processMessage(u *user.User, room *room.Room, msgReq *message.MessageRequest) {
	if err := message.ValidateMessage(msgReq.Content, h.config.MaxMessageLength); err != nil {
		h.sendError(u.Conn, err.Error())
		return
	}

	msg := message.NewMessage(
		uuid.New().String(),
		u.ID,
		u.Username,
		msgReq.Content,
		room.ID,
		"message",
	)

	room.AddMessage(msg)

	h.broadcastMessage(room, msg)
}

func (h *Handlers) broadcastMessage(room *room.Room, msg *message.Message) {
	users := room.GetUsers()

	for _, user := range users {
		if err := user.SendMessage(map[string]interface{}{
			"type": "message",
			"data": msg,
		}); err != nil {
			log.Printf("Ошибка отправки сообщения пользователю %s: %v", user.Username, err)
			room.RemoveUser(user)
			user.Close()
		}
	}
}

func (h *Handlers) broadcastUserJoin(room *room.Room, user *user.User) {
	joinMsg := message.NewMessage(
		uuid.New().String(),
		user.ID,
		user.Username,
		fmt.Sprintf("%s присоединился к чату", user.Username),
		room.ID,
		"join",
	)

	room.AddMessage(joinMsg)
	h.broadcastMessage(room, joinMsg)
}

func (h *Handlers) sendMessageHistory(conn *websocket.Conn, room *room.Room) {
	recentMessages := room.GetRecentMessages(50)

	response := map[string]interface{}{
		"type": "history",
		"data": map[string]interface{}{
			"messages": recentMessages,
			"users":    room.GetUserList(),
		},
	}

	if err := conn.WriteJSON(response); err != nil {
		log.Printf("Ошибка отправки истории сообщений: %v", err)
	}
}

func (h *Handlers) sendError(conn *websocket.Conn, errorMsg string) {
	response := map[string]interface{}{
		"type":    "error",
		"message": errorMsg,
	}

	conn.WriteJSON(response)
	conn.Close()
}
