package phada

import (
	"errors"
	"time"

	"github.com/go-redis/redis"
)

// RedisSessionStore
type RedisSessionStore struct {
	SessionStore
	lastWriteTime time.Time
	client        *redis.Client
}

// NewRedisSessionStore
//
// Creates an inmemory store that uses a concurrent map to store sessions
func NewRedisSessionStore(redisClient *redis.Client) *RedisSessionStore {
	return &RedisSessionStore{
		lastWriteTime: time.Now(),
		client:        redisClient,
	}
}

// PutHop
func (m *RedisSessionStore) PutHop(ussdRequest *UssdRequestSession) error {
	e, ok := m.client.Get(ussdRequest.SessionID)
	if !ok {
		m.client.Set(ussdRequest.SessionID, ussdRequest)
		// return errors.New("Failed to read session data from memory store")
	}

	if e == nil {
		m.client.Set(ussdRequest.SessionID, *ussdRequest)
		return nil
	}
	existing := e.(UssdRequestSession)
	existing.RecordHop(ussdRequest.Text)
	m.client.Set(ussdRequest.SessionID, existing)
	m.lastWriteTime = time.Now()
	return nil
}

// Delete
func (m *RedisSessionStore) Delete(sessionID string) {
	m.client.Del(sessionID)
}

// Get
func (m *RedisSessionStore) Get(sessionID string) (*UssdRequestSession, error) {
	e, ok := m.client.Get(sessionID)
	if !ok {
		return nil, errors.New("Session does not exist in SessionStore")
	}

	session := e.(UssdRequestSession)

	return &session, nil
}
