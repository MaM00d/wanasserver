package chat

import (
	"log/slog"

	api "Server/elapi"
	db "Server/eldb"
)

type ElChat struct {
	ap *api.ElApi
	db *db.Storage
}

// create object of struct apiserver to set the listen addr
func NewElChat(db *db.Storage, elapi *api.ElApi) *ElChat {
	elchat := &ElChat{
		ap: elapi,
		db: db,
	}
	return elchat
}

func (elchat *ElChat) AddRoutes() {
	slog.Info("ElChat Routes")
	elchat.ap.Route("/persoona", elchat.createchat, "POST")
	elchat.ap.Route("/getchats", elchat.getchats, "POST")
}
