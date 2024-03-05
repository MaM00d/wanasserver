package persona

import (
	"database/sql"
	"fmt"
	"log/slog"

	_ "github.com/lib/pq"

	db "Server/eldb"
)

type PersonaStorage interface {
	InsertPersona(*Persona) error
	DeletePersona(int) error
	UpdatePersona(*Persona) error
	GetPersonaById(int) (*Persona, error)
	GetPersonaByEmail(string) (*Persona, error)
}

type personaStore struct {
	db db.Storage
}

func newPersonaStore(db db.Storage) *personaStore {
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
	err := s.db.Exec(query)
	return err
}

func (s *personaStore) dropfunctionid() error {
	query := `
    drop function if exists fn_trig_persona_pk;
    `
	err := s.db.Exec(query)
	return err
}

func (s *personaStore) droptrigid() error {
	query := `
    drop trigger if exists trig_persona_pk on persona;
    `
	err := s.db.Exec(query)
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
	err := s.db.Exec(query)
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
	err := s.db.Exec(query)
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
	err := s.db.Exec(query)
	return err
}

func (s *personaStore) InsertPersona(elpersona *Persona) error {
	query := `insert into Persona 
    (name,userid,createdat)
    values ($1,$2,$3)
    `
	err := s.db.Query(
		query,
		&elpersona.Name,
		&elpersona.UserID,
		&elpersona.CreatedAt)
	if err != nil {
		slog.Error("inserting to database")
		return err
	}

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
	var personas []Persona
	rows := s.db.QueryScan(personas, `select * from Personas where ID = $1`, id)

	if rows == fmt.Errorf("not found") {

		slog.Error("GetPersonasByUserId", "id", id)
		return nil, fmt.Errorf("persona with id [%d] not found", id)
	}
	return personas, nil
}
