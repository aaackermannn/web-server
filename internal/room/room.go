package room

import (
	"sync"
	"time"

	"chat-server/internal/message"
	"chat-server/internal/user"
)

type Room struct {
	ID           string              `json:"id"`
	Name         string              `json:"name"`
	Users        map[*user.User]bool `json:"-"`
	Messages     []*message.Message  `json:"-"`
	CreatedAt    time.Time           `json:"createdAt"`
	LastActivity time.Time           `json:"lastActivity"`
	mutex        sync.RWMutex
}

func NewRoom(id, name string) *Room {
	now := time.Now()
	return &Room{
		ID:           id,
		Name:         name,
		Users:        make(map[*user.User]bool),
		Messages:     make([]*message.Message, 0),
		CreatedAt:    now,
		LastActivity: now,
	}
}

func (r *Room) AddUser(u *user.User) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.Users[u] = true
	r.LastActivity = time.Now()
}

func (r *Room) RemoveUser(u *user.User) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	delete(r.Users, u)
	r.LastActivity = time.Now()
}

func (r *Room) GetUserCount() int {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	return len(r.Users)
}

func (r *Room) GetUsers() []*user.User {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	users := make([]*user.User, 0, len(r.Users))
	for u := range r.Users {
		users = append(users, u)
	}
	return users
}

func (r *Room) AddMessage(msg *message.Message) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.Messages = append(r.Messages, msg)
	r.LastActivity = time.Now()

	if len(r.Messages) > 1000 {
		r.Messages = r.Messages[len(r.Messages)-1000:]
	}
}

func (r *Room) GetRecentMessages(count int) []*message.Message {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	if count <= 0 || count > len(r.Messages) {
		count = len(r.Messages)
	}

	start := len(r.Messages) - count
	if start < 0 {
		start = 0
	}

	return r.Messages[start:]
}

func (r *Room) IsEmpty() bool {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	return len(r.Users) == 0
}

func (r *Room) GetUserList() []map[string]interface{} {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	users := make([]map[string]interface{}, 0, len(r.Users))
	for u := range r.Users {
		users = append(users, map[string]interface{}{
			"id":       u.ID,
			"username": u.Username,
			"joinedAt": u.JoinedAt,
		})
	}
	return users
}
