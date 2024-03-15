package persona

import (
	"Server/user"
	"log/slog"
	"net/http"
)

func (s *ElPersona) getpersonas(w http.ResponseWriter, r *http.Request) error {
	slog.Info("Handling Login")
	eluserid := user.Getidfromheader(r)
	elpersonas, err := s.GetPersonasByUserId(eluserid)
	if err != nil {
		return err
	}
	return s.ap.WriteJSON(w, http.StatusOK, &elpersonas)
}
