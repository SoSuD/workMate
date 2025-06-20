package mapstore

import (
	"github.com/google/uuid"
	"sync"
	"workMate/internal/service/faketask"
	"workMate/internal/store"
)

type Store struct {
	db             map[uuid.UUID]*faketask.Task
	userRepository *TaskRepository
	mu             sync.RWMutex
}

func New() *Store {
	return &Store{
		db: make(map[uuid.UUID]*faketask.Task),
		mu: sync.RWMutex{},
	}
}

func (s *Store) Task() store.TaskRepository {
	if s.userRepository != nil {
		return s.userRepository
	}
	s.userRepository = &TaskRepository{
		store: s,
	}
	return s.userRepository
}
