package phada

import (
	"errors"
	"time"

	"github.com/orcaman/concurrent-map"
)

// InMemorySessionStore
type InMemorySessionStore struct {
	SessionStore
	lastWriteTime time.Time
	data          cmap.ConcurrentMap
}

// NewInMemorySessionStore
//
// Creates an inmemory store that uses a concurrent map to store sessions
func NewInMemorySessionStore() *InMemorySessionStore {
	return &InMemorySessionStore{
		lastWriteTime: time.Now(),
		data:          cmap.New(),
	}
}

// PutHop
func (m *InMemorySessionStore) PutHop(ussdRequest *UssdRequestSession) error {
	e, ok := m.data.Get(ussdRequest.SessionID)
	if !ok {
		m.data.Set(ussdRequest.SessionID, ussdRequest)
		// return errors.New("Failed to read session data from memory store")
	}

	if e == nil {
		m.data.Set(ussdRequest.SessionID, *ussdRequest)
		return nil
	}
	existing := e.(UssdRequestSession)
	existing.RecordHop(ussdRequest.Text)
	m.data.Set(ussdRequest.SessionID, existing)
	m.lastWriteTime = time.Now()
	return nil
}

// Delete
func (m *InMemorySessionStore) Delete(sessionID string) {
	m.data.Remove(sessionID)
}

// Get
func (m *InMemorySessionStore) Get(sessionID string) (*UssdRequestSession, error) {
	e, ok := m.data.Get(sessionID)
	if !ok {
		return nil, errors.New("Session does not exist in SessionStore")
	}

	session := e.(UssdRequestSession)

	return &session, nil
}
