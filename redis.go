package phada

import (
	"encoding/json"
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
	data, err := m.client.Get(ussdRequest.SessionID).Result()
	if err != nil {
		return m.client.Set(ussdRequest.SessionID, ussdRequest.ToJSON(), 0).Err()
	}

	if data == "" {
		return m.client.Set(ussdRequest.SessionID, ussdRequest.ToJSON(), 0).Err()
	}
	var existing *UssdRequestSession
	err = json.Unmarshal([]byte(data), existing)
	if err != nil {
		return err
	}
	existing.RecordHop(ussdRequest.Text)
	err = m.client.Set(ussdRequest.SessionID, existing.ToJSON(), 0).Err()
	if err != nil {
		m.lastWriteTime = time.Now()
	}
	return err
}

// Delete
func (m *RedisSessionStore) Delete(sessionID string) {
	m.client.Del(sessionID)
}

// Get
func (m *RedisSessionStore) Get(sessionID string) (*UssdRequestSession, error) {
	data, err := m.client.Get(sessionID).Result()
	if err != nil {
		return nil, errors.New("Session does not exist in SessionStore")
	}
	var existing *UssdRequestSession
	err = json.Unmarshal([]byte(data), existing)
	if err != nil {
		return nil, err
	}

	return existing, nil
}
