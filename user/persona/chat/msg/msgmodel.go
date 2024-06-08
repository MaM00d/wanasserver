package msg

import (
	"time"
)

type Msg struct {
	ID        int       `db:"id"        json:"id"`
	ChatID    int       `db:"chatid"    json:"chatid"`
	PersonaID int       `db:"personaid" json:"personaid"`
	UserID    int       `db:"userid"    json:"userid"`
	Message   *string   `db:"message"   json:"message"`
	State     bool      `db:"state"     json:"state"`
	CreatedAt time.Time `db:"createdat" json:"createdat"`
}

type MsgView struct {
	Message   string    `db:"message"   json:"message"`
	State     bool      `db:"state"     json:"state"`
	CreatedAt time.Time `db:"createdat" json:"createdat"`
}

func NewMsg(chatid, personaid, userid int, message *string, state bool) *Msg {
	return &Msg{
		ChatID:    chatid,
		PersonaID: personaid,
		UserID:    userid,
		Message:   message,
		State:     state,
		CreatedAt: time.Now().UTC(),
	}
}
