package chat

import (
	"errors"
	"log/slog"

	api "Server/elapi"
	db "Server/eldb"
)

type ElChat struct {
	ap                   *api.ElApi
	db                   *db.Storage
	PersonaDoesnotExsist error
}

// create object of struct apiserver to set the listen addr
func NewElChat(db *db.Storage, elapi *api.ElApi) *ElChat {
	elchat := &ElChat{
		ap:                   elapi,
		db:                   db,
		PersonaDoesnotExsist: errors.New("persona doesn't exsist"),
	}
	return elchat
}

func (elchat *ElChat) AddRoutes() {
	slog.Info("ElChat Routes")
	elchat.ap.Route("/persona/{personaid}", elchat.createchat, "PUT")
	elchat.ap.Route("/persona/{personaid}", elchat.getchats, "GET")
}

func (s *ElChat) InitDb() error {
	if err := s.createChatTabel(); err != nil {
		return err
	}
	if err := s.createpersonafk(); err != nil {
		return err
	}
	if err := s.createuserfk(); err != nil {
		return err
	}

	if err := s.createfunctionid(); err != nil {
		return err
	}
	if err := s.createtriggerid(); err != nil {
		return err
	}

	return nil
}

func (s *ElChat) DropDb() error {
	if err := s.droppersonafk(); err != nil {
		return err
	}
	if err := s.dropuserfk(); err != nil {
		return err
	}
	if err := s.droptrigid(); err != nil {
		return err
	}
	if err := s.dropfunctionid(); err != nil {
		return err
	}
	if err := s.dropChatTabel(); err != nil {
		return err
	}

	return nil
}
