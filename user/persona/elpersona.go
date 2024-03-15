package persona

import (
	"log/slog"

	api "Server/elapi"
	db "Server/eldb"
)

type ElPersona struct {
	ap *api.ElApi
	db *db.Storage
}

// create object of struct apiserver to set the listen addr
func NewElPersona(db *db.Storage, elapi *api.ElApi) *ElPersona {
	elpersona := &ElPersona{
		ap: elapi,
		db: db,
	}
	return elpersona
}

func (elpersona *ElPersona) AddRoutes() {
	slog.Info("ElPersona Routes")
	elpersona.ap.Route("/persona", elpersona.createpersona, "POST")
	elpersona.ap.Route("/persona", elpersona.getpersonas, "GET")
}

func (s *ElPersona) InitDb() error {
	if err := s.createPersonaTabel(); err != nil {
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

func (s *ElPersona) DropDb() error {
	if err := s.dropuserfk(); err != nil {
		return err
	}
	if err := s.droptrigid(); err != nil {
		return err
	}
	if err := s.dropfunctionid(); err != nil {
		return err
	}
	if err := s.dropPersonaTabel(); err != nil {
		return err
	}

	return nil
}
