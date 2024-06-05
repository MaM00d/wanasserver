package main

import (
	"log/slog"
	"net/http"

	ai "Server/elai"
	api "Server/elapi"
	db "Server/eldb"
	user "Server/user"
	persona "Server/user/persona"
	chat "Server/user/persona/chat"
	msg "Server/user/persona/chat/msg"
)

type APIServer struct {
	listenAddr string
	store      *db.Storage
}

func NewApiServer(listenAddr string, store *db.Storage) *APIServer {
	slog.Info("createing server")
	return &APIServer{
		listenAddr: listenAddr,
		store:      store,
	}
}

func (s *APIServer) Run() {
	slog.Info("start")
	ap := api.NewElApi()
	usr := user.NewElUser(s.store, ap)
	pers := persona.NewElPersona(s.store, ap)
	chat := chat.NewElChat(s.store, ap)
	ais := ai.InitAiServer()
	msg := msg.NewElMsg(s.store, ap, ais)
	InitRoutes(usr, pers, chat, msg)
	DropDb(msg, chat, pers, usr)
	InitDb(usr, pers, chat, msg)
	// logging
	slog.Info("JSON API server runngin", "PORT", s.listenAddr)
	// start listening on addresss and sending to router
	http.ListenAndServe(s.listenAddr, ap.GetRouter())
}

type entity interface {
	create() error
	drop() error
}
type object interface {
	InitDb() error
	DropDb() error
	AddRoutes()
}

func DropDb(obj ...object) error {
	for _, o := range obj {
		if err := o.DropDb(); err != nil {
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
