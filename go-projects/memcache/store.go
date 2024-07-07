package memcache

import (
	"sync"
	"time"
)

type Store struct {
	Data              map[string]Data
	DefaultExpiration time.Duration
	mu                sync.RWMutex
}

type Data struct {
	Data           any
	ExpirationTime time.Time
}

func NewStore(expiration time.Duration) *Store {
	return &Store{
		DefaultExpiration: expiration,
		Data:              make(map[string]Data),
	}
}

func (s *Store) Set(key string, data any) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Data[key] = Data{
		Data:           data,
		ExpirationTime: time.Now().Add(s.DefaultExpiration),
	}
}

func (s *Store) Get(key string) (any, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	d, ok := s.Data[key]
	if !ok {
		return nil, false
	}
	return d.Data, true
}

func (s *Store) Delete(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.Data, key)
}

func (s *Store) CleanUp() {
	s.mu.Lock()
	s.Data = make(map[string]Data)
	s.mu.Unlock()
}

func (s *Store) Flush() {
	// TODO: Implement this
	panic("not implemented")
}
