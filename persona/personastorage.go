package persona

import (
	"database/sql"
	"fmt"
	"log/slog"

	_ "github.com/lib/pq"
)

type PersonaStorage interface {
	InsertPersona(*Persona) error
	DeletePersona(int) error
	UpdatePersona(*Persona) error
}

type personaStore struct {
	db *sql.DB
}

func newPersonaStore(db *sql.DB) *personaStore {
	return &personaStore{
		db: db,
	}
}

func (s *personaStore) InitDb() error {
	s.dropPersonaTabel()
	return s.createPersonaTabel()
}

func (s *personaStore) dropPersonaTabel() error {
	query := `
    drop table if exists "Persona";
    `
	_, err := s.db.Exec(query)
	return err
}

func (s *personaStore) createPersonaTabel() error {
	query := `
        CREATE TABLE "Persona" (
            "ID" int   NOT NULL,
            "Name" char(50)   NOT NULL,
            "UserID" int   NOT NULL,
            "CreatedAt" timestamp   NOT NULL,
            CONSTRAINT "pk_Persona" PRIMARY KEY (
                "ID"
             )
        );
    `
	_, err := s.db.Exec(query)
	return err
}

func (s *personaStore) InsertPersona(persona *Persona) error {
	query := `insert into "Persona" 
    ("Name","UserID","CreatedAt")
    values ($1,$2,$3)
    `
	resp, err := s.db.Query(
		query,
		persona.Name,
		persona.UserID,
		persona.CreatedAt,
	)
	if err != nil {
		slog.Error("inserting to database")
		return err
	}

	fmt.Printf("%+v\n", resp)
	return nil
}

func (s *personaStore) DeletePersona(int) error {
	return nil
}

func (s *personaStore) UpdatePersona(*Persona) error {
	return nil
}

// func (s *personaStore) GetPersonaById(int) (*Persona, error) {
// 	return nil, nil
// }

func scanIntoAccount(rows *sql.Rows) (*Persona, error) {
	elpersona := new(Persona)
	err := rows.Scan(
		&elpersona.Name,
		&elpersona.UserID,
		&elpersona.CreatedAt)

	return elpersona, err
}

func (s *personaStore) GetPersonasByUserID(userid int) (*Persona, error) {
	elpersona := new(Persona)
	personas, err := s.db.Query("select * from Persona where UserID = $1", userid)

	if err == sql.ErrNoRows {

		slog.Info("no persona found with this userid", "email", userid)
		return nil, fmt.Errorf("persona with userid [%d] not found", userid)
	}

	for personas.Next() {
		return scanIntoAccount(personas)
	}

	return elpersona, nil
}
