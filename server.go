package main

import (
	"log/slog"
	"net/http"

	api "Server/elapi"
	db "Server/eldb"
	user "Server/user"
)

// persona "Server/persona"
type APIServer struct {
	listenAddr string
	store      *db.Storage
}

func NewApiServer(listenAddr string, store *db.Storage) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		store:      store,
	}
}

func (s *APIServer) Run() {
	ap := api.NewElApi()
	usr := user.NewElUser(s.store, ap)
	// pers := persona.NewElPersona(s.store, ap)
	InitRoutes(usr)
	InitDb(usr)
	// DropDb(pers, usr)
	// chat.NewElMsg(s.store.db, ap)
	// logging
	slog.Info("JSON API server runngin", "PORT", s.listenAddr)
	// start listening on addresss and sending to router
	http.ListenAndServe(s.listenAddr, ap.GetRouter())
}

type object interface {
	InitDb() error
	DropDb() error
	AddRoutes()
}

func DropDb(obj ...object) error {
	for _, o := range obj {
		if err := o.InitDb(); err != nil {
			return err
		}
	}
	return nil
}

func InitDb(obj ...object) error {
	for _, o := range obj {
		if err := o.InitDb(); err != nil {
			return err
		}
	}
	return nil
}

func InitRoutes(obj ...object) {
	for _, o := range obj {
		o.AddRoutes()
	}
}
