package persona

import (
	"log/slog"

	api "Server/elapi"
	db "Server/eldb"
)

type ElPersona struct {
	ap *api.ElApi
	db *db.Storage
}

// create object of struct apiserver to set the listen addr
func NewElPersona(db *db.Storage, elapi *api.ElApi) *ElPersona {
	elpersona := &ElPersona{
		ap: elapi,
		db: db,
	}
	return elpersona
}

func (elpersona *ElPersona) AddRoutes() {
	slog.Info("ElPersona Routes")
	elpersona.ap.Route("/persoona", elpersona.createpersona, "POST")
	elpersona.ap.Route("/getpersonas", elpersona.getpersonas, "POST")
}
