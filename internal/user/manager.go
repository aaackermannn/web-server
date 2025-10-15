package user

import (
	"sync"
	"time"
)

type Manager struct {
	users map[string]*User
	mutex sync.RWMutex
}

func NewManager() *Manager {
	manager := &Manager{
		users: make(map[string]*User),
	}

	go manager.checkUserActivity()

	return manager
}

func (m *Manager) AddUser(user *User) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.users[user.ID] = user
}

func (m *Manager) RemoveUser(userID string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	delete(m.users, userID)
}

func (m *Manager) GetUser(userID string) (*User, bool) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	user, exists := m.users[userID]
	return user, exists
}

func (m *Manager) GetAllUsers() []*User {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	users := make([]*User, 0, len(m.users))
	for _, user := range m.users {
		users = append(users, user)
	}
	return users
}

func (m *Manager) GetOnlineUserCount() int {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return len(m.users)
}

func (m *Manager) checkUserActivity() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		m.mutex.Lock()
		for userID, user := range m.users {
			if !user.IsAlive() {
				user.Close()
				delete(m.users, userID)
			}
		}
		m.mutex.Unlock()
	}
}
