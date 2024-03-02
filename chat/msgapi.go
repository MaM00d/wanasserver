package msg

import (
	"database/sql"

	api "Server/elapi"
)

type ElChat struct {
	ap    *api.ElApi
	store ChatStorage
}

// create object of struct apiserver to set the listen addr
func NewElMsg(db *sql.DB, elapi *api.ElApi) *ElChat {
	store := newChatStore(db)
	elChat := &ElChat{
		ap:    elapi,
		store: store,
	}
	elChat.addroutes()
	return elChat
}

func (elmsg *ElChat) addroutes() {
	elmsg.ap.Route("/msg", ElChat.Send, "POST")
}
