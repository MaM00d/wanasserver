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
	GetPersonaById(int) (*Persona, error)
	GetPersonaByEmail(string) (*Persona, error)
}

type personaStore struct {
	db *sql.DB
}

func newPersonaStore(db *sql.DB) *personaStore {
	elstore := &personaStore{
		db: db,
	}
	elstore.InitDb()
	return elstore
}

func (s *personaStore) InitDb() error {
	s.droptrigid()
	s.dropfunctionid()
	s.dropPersonaTabel()
	s.createPersonaTabel()
	s.createfunctionid()
	return s.createtriggerid()
}

func (s *personaStore) dropPersonaTabel() error {
	query := `
    drop table if exists Persona;
    `
	_, err := s.db.Exec(query)
	return err
}

func (s *personaStore) dropfunctionid() error {
	query := `
    drop function if exists fn_trig_persona_pk;
    `
	_, err := s.db.Exec(query)
	return err
}

func (s *personaStore) droptrigid() error {
	query := `
    drop trigger if exists trig_persona_pk on persona;
    `
	_, err := s.db.Exec(query)
	return err
}

func (s *personaStore) createPersonaTabel() error {
	query := `
            CREATE TABLE persona (
                id int   ,
                name char(50)   NOT NULL,
                userid int   NOT NULL,
                createdat timestamp   NOT NULL,
                CONSTRAINT pk_persona PRIMARY KEY (
                    id,useriD
                 )
            );
    `
	_, err := s.db.Exec(query)
	return err
}

func (s *personaStore) createtriggerid() error {
	query := `

            CREATE TRIGGER trig_persona_pk
              BEFORE insert 
              ON persona
              FOR EACH ROW
              EXECUTE PROCEDURE fn_trig_persona_pk();
    `
	_, err := s.db.Exec(query)
	return err
}

func (s *personaStore) createfunctionid() error {
	query := `
            CREATE OR REPLACE FUNCTION "fn_trig_persona_pk"()
              RETURNS "pg_catalog"."trigger" AS $BODY$ 
            begin
            new.id = (select count(*)+1 from persona where userid=new.userid);
            return NEW;
            end;
            $BODY$
              LANGUAGE plpgsql VOLATILE
              COST 100;

    `
	_, err := s.db.Exec(query)
	return err
}

func (s *personaStore) InsertPersona(elpersona *Persona) error {
	query := `insert into Persona 
    (name,userid,createdat)
    values ($1,$2,$3)
    `
	resp, err := s.db.Query(
		query,
		&elpersona.Name,
		&elpersona.UserID,
		&elpersona.CreatedAt)
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
		&elpersona.ID,
		&elpersona.Name,
		&elpersona.UserID,
		&elpersona.CreatedAt)

	return elpersona, err
}

func (s *personaStore) GetPersonasByUserId(id int) ([]Persona, error) {
	rows, err := s.db.Query(`select * from Personas where ID = $1`, id)

	if err == sql.ErrNoRows {

		slog.Info("no persona found with this email", "email", id)
		return nil, fmt.Errorf("persona with email [%d] not found", id)
	}
	var personas []Persona
	for rows.Next() {
		var persona Persona
		if err := rows.Scan(&persona.ID, &persona.Name, &persona.UserID,
			&persona.CreatedAt); err != nil {
			return personas, err
		}
		personas = append(personas, persona)
	}

	return personas, nil
}
