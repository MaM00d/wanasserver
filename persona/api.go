package persona

import (
	"database/sql"

	api "Server/elapi"
)

type ElPersona struct {
	ap    *api.ElApi
	store *personaStore
}

// create object of struct apiserver to set the listen addr
func NewElPersona(db *sql.DB, elapi *api.ElApi) *ElPersona {
	store := newPersonaStore(db)
	elpersona := &ElPersona{
		ap:    elapi,
		store: store,
	}
	elpersona.addroutes()
	return elpersona
}

func (elpersona *ElPersona) addroutes() {
	elpersona.ap.Route("/persona", elpersona.createpersona, "POST")
	elpersona.ap.Route("/getpersonas", elpersona.getpersonas, "POST")
}
