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
	elpersonaid, err := strconv.Atoi(s.ap.GetFromVars(r, "personaid"))
	elchats, err := s.GetChatsByUserId(elpersonaid, eluserid)
	if err != nil {
		return err
	}
	return s.ap.WriteJSON(w, http.StatusOK, elchats)
}
