package ratelimit

import "time"

// Store is the limit storage
type Store interface {
	Incr(key string, ttl time.Duration) (int64, error)
}

type memoryStore struct {
	data map[string]int64
}

// NewMemoryStore creates new memory limit storage
func NewMemoryStore() Store {
	s := memoryStore{}
	s.data = make(map[string]int64)
	return &s
}

func (s *memoryStore) Incr(key string, ttl time.Duration) (int64, error) {
	s.data[key]++
	return s.data[key], nil
}
