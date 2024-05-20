package persona

import (
	"Server/user"
	"log/slog"
	"net/http"
)

func (s *ElPersona) getpersonas(w http.ResponseWriter, r *http.Request) error {
	slog.Info("Handling Login")
	eluserid := user.Getidfromheader(r)
	if eluserid < 0 {
		s.ap.WriteError(w, http.StatusUnauthorized, "invalid token")
	}
	elpersonas, err := s.GetPersonasByUserId(eluserid)
	if err != nil {
		if err == s.db.NotFound {
			return s.ap.WriteError(w, http.StatusOK, "")
		}

		return err

	}

	return s.ap.WriteJSON(w, http.StatusOK, &elpersonas)
}
