package store

type Store interface {
	Task() TaskRepository
}
