package main

import (
	"log/slog"
	"net/http"

	api "Server/elapi"
	user "Server/user"
)

type APIServer struct {
	listenAddr string
	store      Storage
}

func NewApiServer(listenAddr string, store *Storage) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		store:      *store,
	}
}

func (s *APIServer) Run() {
	ap := api.NewElApi()
	user.NewElUser(s.store.db, ap)
	// logging
	slog.Info("JSON API server runngin", "PORT", s.listenAddr)
	// start listening on addresss and sending to router
	http.ListenAndServe(s.listenAddr, ap.GetRouter())
}
