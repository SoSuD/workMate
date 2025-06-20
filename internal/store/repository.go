package store

import (
	"github.com/google/uuid"
	"workMate/internal/service/faketask"
)

type TaskRepository interface {
	Create(name string) (*faketask.Task, error)
	Get(id uuid.UUID) (*faketask.Task, error)
	Delete(id uuid.UUID) error
	Finish(id uuid.UUID, status string, values string) error
}
