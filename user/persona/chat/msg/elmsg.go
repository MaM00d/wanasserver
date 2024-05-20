package msg

import (
	"errors"
	"log/slog"

	ai "Server/elai"
	api "Server/elapi"
	db "Server/eldb"
)

type ElMsg struct {
	ap                   *api.ElApi
	db                   *db.Storage
	PersonaDoesnotExsist error
	ais                  *ai.Aiserver
}

// create object of struct apiserver to set the listen addr
func NewElMsg(db *db.Storage, elapi *api.ElApi) *ElMsg {
	ais := ai.InitAiServer("localhost:12345")

	elmsg := &ElMsg{
		ap:                   elapi,
		db:                   db,
		PersonaDoesnotExsist: errors.New("persona doesn't exsist"),
		ais:                  ais,
	}
	return elmsg
}

func (elmsg *ElMsg) AddRoutes() {
	slog.Info("ElMsg Routes")
	elmsg.ap.Route("/persona/{personaid}/chat/{chatid}", elmsg.sendmsg, "POST")
	elmsg.ap.Route("/persona/{personaid}/chat/{chatid}", elmsg.getmsgs, "GET")
}

func (s *ElMsg) InitDb() error {
	if err := s.createMsgTabel(); err != nil {
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

func (s *ElMsg) DropDb() error {
	if err := s.droptrigid(); err != nil {
		return err
	}
	if err := s.dropfunctionid(); err != nil {
		return err
	}
	if err := s.dropMsgTabel(); err != nil {
		return err
	}

	return nil
}
