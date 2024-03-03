package chat

import (
	"database/sql"

	api "Server/elapi"
)

type ElChat struct {
	ap    *api.ElApi
	store *chatStore
}

// create object of struct apiserver to set the listen addr
func NewElChat(db *sql.DB, elapi *api.ElApi) *ElChat {
	store := newChatStore(db)
	elchat := &ElChat{
		ap:    elapi,
		store: store,
	}
	elchat.addroutes()
	return elchat
}

func (elchat *ElChat) addroutes() {
	elchat.ap.Route("/sendmsg", elchat.createchat, "POST")
	elchat.ap.Route("/getchat", elchat.getchats, "POST")
}
