package persona

import (
	"log/slog"

	api "Server/elapi"
	db "Server/eldb"
)

type ElPersona struct {
	ap    *api.ElApi
	store *personaStore
}

// create object of struct apiserver to set the listen addr
func NewElPersona(db db.Storage, elapi *api.ElApi) *ElPersona {
	store := newPersonaStore(db)
	elpersona := &ElPersona{
		ap:    elapi,
		store: store,
	}
	elpersona.addroutes()
	return elpersona
}

func (elpersona *ElPersona) addroutes() {
	slog.Info("ElPersona Routes")
	elpersona.ap.Route("/persoona", elpersona.createpersona, "POST")
	elpersona.ap.Route("/getpersonas", elpersona.getpersonas, "POST")
}
