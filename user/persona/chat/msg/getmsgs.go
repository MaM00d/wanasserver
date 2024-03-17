package msg

import (
	"Server/user"
	"log/slog"
	"net/http"
	"strconv"
)

func (s *ElMsg) getmsgs(w http.ResponseWriter, r *http.Request) error {
	slog.Info("Handling Login")

	eluserid := user.Getidfromheader(r)
	if eluserid < 0 {
		s.ap.WriteError(w, http.StatusUnauthorized, "invalid token")
	}

	elchatid, err := strconv.Atoi(s.ap.GetFromVars(r, "chatid"))
	elpersonaid, err := strconv.Atoi(s.ap.GetFromVars(r, "personaid"))
	elmsgs, err := s.GetMsgs(elchatid, elpersonaid, eluserid)
	if err != nil {
		s.ap.WriteError(w, http.StatusNotFound, "not found url")
	}

	return s.ap.WriteJSON(w, http.StatusOK, elmsgs)
}
