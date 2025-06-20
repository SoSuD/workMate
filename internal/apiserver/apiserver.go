package apiserver

import (
	"fmt"
	"log"
	"workMate/internal/store"
	"workMate/internal/store/mapstore"
)

func Start(config Config) {
	srv := NewServer(NewDB(), config)
	srv.configureRouter()
	if err := srv.router.Run(fmt.Sprintf(":%s", config.Server.Port)); err != nil {
		log.Fatal(err)
	}
}

func NewDB() store.Store {
	return mapstore.New()
}
