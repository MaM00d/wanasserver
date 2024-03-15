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
	elchatid, err := strconv.Atoi(s.ap.GetFromVars(r, "chatid"))
	elpersonaid, err := strconv.Atoi(s.ap.GetFromVars(r, "personaid"))
	elmsgs, err := s.GetMsgs(elchatid, elpersonaid, eluserid)
	if err != nil {
		return err
	}

	return s.ap.WriteJSON(w, http.StatusOK, elmsgs)
}
