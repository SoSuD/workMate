package mapstore

import (
	"github.com/google/uuid"
	"time"
	"workMate/internal/service/faketask"
	"workMate/internal/store"
)

type TaskRepository struct {
	store *Store
}

func (u *TaskRepository) Create(name string) (*faketask.Task, error) {
	id := uuid.New()
	t := &faketask.Task{
		ID:        id,
		CreatedAt: time.Now(),
		Status:    "CREATED",
		Name:      name,
	}

	u.store.mu.Lock()
	defer u.store.mu.Unlock()

	u.store.db[id] = t
	return t, nil
}

func (u *TaskRepository) Delete(id uuid.UUID) error {
	u.store.mu.RLock()
	_, ok := u.store.db[id]
	u.store.mu.RUnlock()

	if !ok {
		return store.ErrTaskNotFound
	}

	u.store.mu.Lock()
	defer u.store.mu.Unlock()
	delete(u.store.db, id)
	return nil
}

func (u *TaskRepository) Get(id uuid.UUID) (*faketask.Task, error) {
	u.store.mu.RLock()
	defer u.store.mu.RUnlock()
	t, ok := u.store.db[id]
	if !ok {
		return nil, store.ErrTaskNotFound
	}
	return t, nil
}

func (u *TaskRepository) Finish(id uuid.UUID, status string, value string) error {
	u.store.mu.RLock()
	_, ok := u.store.db[id]
	u.store.mu.RUnlock()
	if !ok {
		return store.ErrTaskNotFound
	}

	u.store.mu.Lock()
	defer u.store.mu.Unlock()

	u.store.db[id].Status = status
	u.store.db[id].FinishedAt = time.Now()
	u.store.db[id].Result = value
	return nil

}
