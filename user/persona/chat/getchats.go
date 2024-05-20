package chat

import (
	"Server/user"
	"log/slog"
	"net/http"
	"strconv"
)

func (s *ElChat) getchats(w http.ResponseWriter, r *http.Request) error {
	slog.Info("Handling getchats")
	eluserid := user.Getidfromheader(r)
	if eluserid < 0 {
		s.ap.WriteError(w, http.StatusUnauthorized, "invalid token")
	}
	personaid := s.ap.GetFromVars(r, "personaid")

	elpersonaid, err := strconv.Atoi(personaid)
	elchats, err := s.GetChatsByUserId(elpersonaid, eluserid)
	if err != nil {
		if err == s.db.NotFound {
			return s.ap.WriteJSON(w, http.StatusOK, "")
		}
		if err == s.PersonaDoesnotExsist {
			return s.ap.WriteError(w, http.StatusNotFound, "Persona doesn't exsist")
		}
		return err
	}
	return s.ap.WriteJSON(w, http.StatusOK, elchats)
}
