package msg

import (
	"Server/user"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
)

type MsgRequest struct {
	Message string `db:"message" json:"message"`
}
type MsgResponse struct {
	Message string `db:"message" json:"message"`
}

func (s *ElMsg) sendmsg(w http.ResponseWriter, r *http.Request) error {
	// decode json from request
	slog.Info("Handling Create Msg")
	msgReq := new(MsgRequest)
	if err := json.NewDecoder(r.Body).Decode(msgReq); err != nil {
		slog.Error("decoding request body", "Model", "Msg")
		return err
	}

	eluserid := user.Getidfromheader(r)
	if eluserid < 0 {
		s.ap.WriteError(w, http.StatusUnauthorized, "invalid token")
	}

	elchatid, err := strconv.Atoi(s.ap.GetFromVars(r, "chatid"))
	elpersonaid, err := strconv.Atoi(s.ap.GetFromVars(r, "personaid"))
	msg := NewMsg(
		elchatid,
		elpersonaid,
		eluserid,
		msgReq.Message,
	)
	if err != nil {
		return err
	}

	if err := s.InsertMsg(msg); err != nil {
		return err
	}

	resp := MsgResponse{
		Message: "hello from ai",
	}
	slog.Info("Successfully Registered")
	return s.ap.WriteJSON(w, http.StatusOK, resp)
}
