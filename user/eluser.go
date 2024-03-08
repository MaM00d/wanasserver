package user

import (
	api "Server/elapi"
	db "Server/eldb"
)

type ElUser struct {
	ap *api.ElApi
	db *db.Storage
}

// create object of struct apiserver to set the listen addr
func NewElUser(db *db.Storage, elapi *api.ElApi) *ElUser {
	eluser := &ElUser{
		ap: elapi,
		db: db,
	}

	return eluser
}

func (eluser *ElUser) AddRoutes() {
	eluser.ap.Route("/login", eluser.Login, "POST")
	eluser.ap.Route("/register", eluser.Register, "POST")
	eluser.ap.Route("/user/{id}", eluser.JWTAuthMiddleware(eluser.getUserByEmailFromVars), "GET")
}
