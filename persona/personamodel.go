package persona

import (
	"time"
)

type Persona struct {
	Name      string    `json:"personaname"`
	UserID    int       `json:"userid"`
	CreatedAt time.Time `json:"CreatedAt"`
}

type PersonaView struct {
	Name      string    `json:"personaname"`
	UserID    int       `json:"userid"`
	CreatedAt time.Time `json:"CreatedAt"`
}

func NewPersona(name string, userid int) (*Persona, error) {
	return &Persona{
		Name:      name,
		UserID:    userid,
		CreatedAt: time.Now().UTC(),
	}, nil
}
