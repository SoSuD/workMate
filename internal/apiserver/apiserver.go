package apiserver

import (
	"workMate/internal/store"
	"workMate/internal/store/mapstore"
)

func Start(config Config) {
	srv := NewServer(NewDB(), config)
	srv.configureRouter()
	srv.router.Run(":8080")
}

func NewDB() store.Store {
	return mapstore.New()
}
