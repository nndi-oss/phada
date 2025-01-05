package phada

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	redis "github.com/redis/go-redis/v9"
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
	ctx := context.Background()
	data, err := m.client.Get(ctx, ussdRequest.SessionID).Result()
	if err != nil {
		return m.client.Set(ctx, ussdRequest.SessionID, ussdRequest.ToJSON(), 0).Err()
	}

	if data == "" {
		return m.client.Set(ctx, ussdRequest.SessionID, ussdRequest.ToJSON(), 0).Err()
	}
	var existing *UssdRequestSession
	err = json.Unmarshal([]byte(data), existing)
	if err != nil {
		return err
	}
	existing.RecordHop(ussdRequest.Text)
	err = m.client.Set(ctx, ussdRequest.SessionID, existing.ToJSON(), 0).Err()
	if err != nil {
		m.lastWriteTime = time.Now()
	}
	return err
}

// Delete
func (m *RedisSessionStore) Delete(sessionID string) {
	m.client.Del(context.Background(), sessionID)
}

// Get
func (m *RedisSessionStore) Get(sessionID string) (*UssdRequestSession, error) {
	data, err := m.client.Get(context.Background(), sessionID).Result()
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
