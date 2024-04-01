package mocks

import (
	"context"
	"sync"
)

type MockStore struct {
	data map[string]bool
	mu   sync.RWMutex
}

func NewMockStore() *MockStore {
	return &MockStore{
		data: make(map[string]bool),
	}
}

func (m *MockStore) Set(_ context.Context, key string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.data[key] = true
	return nil
}

func (m *MockStore) Exists(_ context.Context, key string) (bool, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	_, exists := m.data[key]
	return exists, nil
}
