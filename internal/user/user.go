package user

import (
	"time"

	"github.com/gorilla/websocket"
)

type User struct {
	ID       string          `json:"id"`
	Username string          `json:"username"`
	Conn     *websocket.Conn `json:"-"`
	JoinedAt time.Time       `json:"joinedAt"`
	LastPing time.Time       `json:"-"`
}

func NewUser(id, username string, conn *websocket.Conn) *User {
	return &User{
		ID:       id,
		Username: username,
		Conn:     conn,
		JoinedAt: time.Now(),
		LastPing: time.Now(),
	}
}

func (u *User) SendMessage(message interface{}) error {
	return u.Conn.WriteJSON(message)
}

func (u *User) Close() error {
	return u.Conn.Close()
}

func (u *User) UpdatePing() {
	u.LastPing = time.Now()
}

func (u *User) IsAlive() bool {
	return time.Since(u.LastPing) < 2*time.Minute
}
