package user

import (
	api "Server/elapi"
	db "Server/eldb"
)

type ElUser struct {
	ap    *api.ElApi
	store *userStore
}

// create object of struct apiserver to set the listen addr
func NewElUser(db db.Storage, elapi *api.ElApi) *ElUser {
	store := newUserStore(db)
	eluser := &ElUser{
		ap:    elapi,
		store: store,
	}
	eluser.addroutes()
	return eluser
}

func (eluser *ElUser) addroutes() {
	eluser.ap.Route("/login", eluser.Login, "POST")
	eluser.ap.Route("/register", eluser.Register, "POST")
	eluser.ap.Route("/user/{id}", eluser.JWTAuthMiddleware, "GET")
}
