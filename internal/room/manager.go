package room

import (
	"sync"
	"time"
)

type Manager struct {
	rooms map[string]*Room
	mutex sync.RWMutex
}

func NewManager() *Manager {
	manager := &Manager{
		rooms: make(map[string]*Room),
	}

	go manager.cleanupEmptyRooms()

	return manager
}

func (m *Manager) GetOrCreateRoom(roomID, roomName string) *Room {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if room, exists := m.rooms[roomID]; exists {
		return room
	}

	room := NewRoom(roomID, roomName)
	m.rooms[roomID] = room
	return room
}

func (m *Manager) GetRoom(roomID string) (*Room, bool) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	room, exists := m.rooms[roomID]
	return room, exists
}

func (m *Manager) DeleteRoom(roomID string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	delete(m.rooms, roomID)
}

func (m *Manager) GetAllRooms() []*Room {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	rooms := make([]*Room, 0, len(m.rooms))
	for _, room := range m.rooms {
		rooms = append(rooms, room)
	}
	return rooms
}

func (m *Manager) GetRoomsInfo() []map[string]interface{} {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	roomsInfo := make([]map[string]interface{}, 0, len(m.rooms))
	for _, room := range m.rooms {
		roomsInfo = append(roomsInfo, map[string]interface{}{
			"id":           room.ID,
			"name":         room.Name,
			"userCount":    room.GetUserCount(),
			"createdAt":    room.CreatedAt,
			"lastActivity": room.LastActivity,
		})
	}
	return roomsInfo
}

func (m *Manager) cleanupEmptyRooms() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		m.mutex.Lock()
		for roomID, room := range m.rooms {
			if room.IsEmpty() {
				delete(m.rooms, roomID)
			}
		}
		m.mutex.Unlock()
	}
}
