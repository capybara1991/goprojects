package store

import (
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Task struct {
	ID string
	Title string
	Description string
	Completed bool
	CreatedAt time.Time
	CompletedAt *time.Time
}

var (
	errNotFound = errors.New("not found")
	errInvalid = errors.New("invalid")
)

type Memory struct {
	mu sync.RWMutex
	t map[string]*Task
	order []string
}

func NewMemory() *Memory {
	return &Memory{t: map[string]*Task{}, order: []string{}}
}

func (m *Memory) CreateTask(title, description string) (Task, error) {
	if title == "" || description == "" { return Task{}, errInvalid }
	m.mu.Lock(); defer m.mu.Unlock()
	id := uuid.NewString()
	now := time.Now()
	t := &Task{ID: id, Title: title, Description: description, Completed: false, CreatedAt: now}
	m.t[id] = t
	m.order = append(m.order, id)
	return *t, nil
}

func (m *Memory) GetTask(id string) (Task, error) {
	m.mu.RLock(); defer m.mu.RUnlock()
	t, ok := m.t[id]
	if !ok { return Task{}, errNotFound }
	cp := *t
	return cp, nil
}

func (m *Memory) SetCompleted(id string, complete bool) (Task, error) {
	m.mu.Lock(); defer m.mu.Unlock()
	t, ok := m.t[id]
	if !ok { return Task{}, errNotFound }
	if complete && !t.Completed {
		now := time.Now()
		t.Completed = true
		t.CompletedAt = &now
	}
	if !complete && t.Completed {
		t.Completed = false
		t.CompletedAt = nil
	}
	cp := *t
	return cp, nil
}

func (m *Memory) ListTasks(onlyUncompleted bool, limit, offset int) []Task {
	m.mu.RLock(); defer m.mu.RUnlock()
	res := make([]Task, 0)
	for _, id := range m.order {
		t := m.t[id]
		if onlyUncompleted && t.Completed { continue }
		res = append(res, *t)
	}
	if offset > len(res) { return []Task{} }
	end := offset + limit
	if limit <= 0 || end > len(res) { end = len(res) }
	out := make([]Task, end-offset)
	copy(out, res[offset:end])
	return out
}

func ToUnix(t *time.Time) int64 {
	if t == nil { return 0 }
	return t.Unix()
}
