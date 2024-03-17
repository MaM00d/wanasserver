package user

import (
	api "Server/elapi"
	db "Server/eldb"
	"log/slog"
)

type ElUser struct {
	ap *api.ElApi
	db *db.Storage
}

// create object of struct apiserver to set the listen addr
func NewElUser(db *db.Storage, elapi *api.ElApi) *ElUser {
	slog.Info("creating api")
	eluser := &ElUser{
		ap: elapi,
		db: db,
	}

	return eluser
}

func (eluser *ElUser) AddRoutes() {
	eluser.ap.Route("/login", eluser.Login, "POST")
	eluser.ap.Route("/register", eluser.Register, "POST")
	eluser.ap.Route("/user/", eluser.getUser, "GET")
	eluser.ap.GetRouter().Use(eluser.AuthMiddleware)
}

func (s *ElUser) InitDb() error {
	if err := s.createUserTable(); err != nil {
		return err
	}

	return nil
}

func (s *ElUser) DropDb() error {
	if err := s.dropUserTabel(); err != nil {
		return err
	}

	return nil
}
