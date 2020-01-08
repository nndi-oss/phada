package phada

import (
	"errors"
	"time"

	"github.com/dgraph-io/ristretto"
)

// RistrettoSessionStore
type RistrettoSessionStore struct {
	SessionStore
	lastWriteTime time.Time
	cache         *ristretto.Cache
}

// NewRistrettoSessionStore
//
// Creates an inmemory store that uses a concurrent map to store sessions
func NewRistrettoSessionStore() *RistrettoSessionStore {
	cache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,     // number of keys to track frequency of (10M).
		MaxCost:     1 << 30, // maximum cost of cache (1GB).
		BufferItems: 64,      // number of keys per Get buffer.
	})
	if err != nil {
		panic(err)
	}
	return &RistrettoSessionStore{
		lastWriteTime: time.Now(),
		cache:         cache,
	}
}

// PutHop
func (m *RistrettoSessionStore) PutHop(ussdRequest *UssdRequestSession) error {
	existing := ussdRequest

	var written bool
	data, found := m.cache.Get(ussdRequest.SessionID)
	if found {
		written = m.cache.Set(ussdRequest.SessionID, ussdRequest, 1)
		if !written {
			return errors.New("Ristretto: Failed to write Session to store")
		}
		existing = data.(*UssdRequestSession)
		existing.RecordHop(ussdRequest.Text)
	}

	written = m.cache.Set(ussdRequest.SessionID, existing, 1)

	if written {
		m.lastWriteTime = time.Now()
	}
	return nil
}

// Delete
func (m *RistrettoSessionStore) Delete(sessionID string) {
	m.cache.Del(sessionID)
}

// Get
func (m *RistrettoSessionStore) Get(sessionID string) (*UssdRequestSession, error) {
	data, found := m.cache.Get(sessionID)
	if !found {
		return nil, errors.New("Ristretto: Session does not exist in SessionStore")
	}

	existing := data.(*UssdRequestSession)

	return existing, nil
}
