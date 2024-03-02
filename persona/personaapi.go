package persona

import (
	"database/sql"

	api "Server/elapi"
)

type ElPersona struct {
	ap    *api.ElApi
	store PersonaStorage
}

// create object of struct apiserver to set the listen addr
func NewElPersona(db *sql.DB, elapi *api.ElApi) *ElPersona {
	store := newPersonaStore(db)
	elPersona := &ElPersona{
		ap:    elapi,
		store: store,
	}
	elPersona.addroutes()
	return elPersona
}

func (elpersona *ElPersona) addroutes() {
	elpersona.ap.Route("/msg", elpersona.getPersonas, "POST")
}
