package msg

import (
	"time"
)

type Msg struct {
	text      string    `json:"text"`
	email     string    `json:"email"`
	personaid string    `json:"persona"`
	CreatedAt time.Time `json:"CreatedAt"`
}

type MsgView struct {
	text      string    `json:"text"`
	email     string    `json:"email"`
	CreatedAt time.Time `json:"CreatedAt"`
}

func NewMsg(text, email, personaid string) (*Msg, error) {
	return &Msg{
		text:      text,
		email:     email,
		personaid: personaid,
		CreatedAt: time.Now().UTC(),
	}, nil
}
