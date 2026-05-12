package match

import (
	"errors"
	"sync"
)

var ErrDuplicate = errors.New("match already exists")

type Registry struct {
	mu sync.RWMutex
	m  map[string]*Match
}

func NewRegistry() *Registry { return &Registry{m: make(map[string]*Match)} }

func (r *Registry) Add(m *Match) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.m[m.ID()]; ok {
		return ErrDuplicate
	}
	r.m[m.ID()] = m
	return nil
}

func (r *Registry) Get(id string) (*Match, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	m, ok := r.m[id]
	return m, ok
}

func (r *Registry) Remove(id string) {
	r.mu.Lock()
	delete(r.m, id)
	r.mu.Unlock()
}
