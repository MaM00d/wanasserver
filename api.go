package main

import (
	"log/slog"
	"net/http"

	api "Server/elapi"
	db "Server/eldb"
	persona "Server/persona"
	user "Server/user"
)

type APIServer struct {
	listenAddr string
	store      db.Storage
}

func NewApiServer(listenAddr string, store *db.Storage) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		store:      *store,
	}
}

func (s *APIServer) Run() {
	ap := api.NewElApi()
	user.NewElUser(s.store, ap)
	persona.NewElPersona(s.store, ap)
	// chat.NewElMsg(s.store.db, ap)
	// logging
	slog.Info("JSON API server runngin", "PORT", s.listenAddr)
	// start listening on addresss and sending to router
	http.ListenAndServe(s.listenAddr, ap.GetRouter())
}
