package msg

import (
	"log/slog"
	"net/http"
	"strconv"

	"Server/user"
)

func (s *ElMsg) getmsgs(w http.ResponseWriter, r *http.Request) error {
	slog.Info("Handling Login")

	eluserid := user.Getidfromheader(r)
	if eluserid < 0 {
		s.ap.WriteError(w, http.StatusUnauthorized, "invalid token")
	}

	elchatid, err := strconv.Atoi(s.ap.GetFromVars(r, "chatid"))
	elpersonaid, err := strconv.Atoi(s.ap.GetFromVars(r, "personaid"))
	elmsgs, err := s.GetMsgs(eluserid, elpersonaid, elchatid)
	if err != nil {
		if err == s.db.NotFound {
			return s.ap.WriteJSON(w, http.StatusOK, "")
		}
		if err == s.PersonaDoesnotExsist {
			return s.ap.WriteError(w, http.StatusNotFound, "Persona doesn't exsist")
		}
		return err
	}
	return s.ap.WriteJSON(w, http.StatusOK, elmsgs)
}
