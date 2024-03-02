package user

import (
	"database/sql"
	"encoding/json"
	"log/slog"
	"net/http"

	api "Server/elapi"
)

type ElUser struct {
	ap    *api.ElApi
	store UserStorage
}

// create object of struct apiserver to set the listen addr
func NewElUser(db *sql.DB, elapi *api.ElApi) *ElUser {
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
	eluser.ap.Route("/user/{email}", eluser.JWTAuthMiddleware, "GET")
}
