package persona

import (
	"database/sql"
	"fmt"
	"log/slog"

	_ "github.com/lib/pq"
)

func (s *ElPersona) createPersonaTabel() error {
	query := `
            CREATE TABLE if not exists persona (
                id int unique  ,
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

func (s *ElPersona) dropPersonaTabel() error {
	query := `
    drop table if exists Persona;
    `
	err := s.db.Exec(query)
	return err
}

func (s *ElPersona) createuserfk() error {
	query := `
    ALTER TABLE persona ADD CONSTRAINT fk_persona_useriD FOREIGN KEY(useriD)
    REFERENCES users (id);
    `
	err := s.db.Exec(query)
	return err
}

func (s *ElPersona) dropuserfk() error {
	query := `
    ALTER TABLE persona
    drop CONSTRAINT fk_persona_userid;
    `
	err := s.db.Exec(query)
	return err
}

func (s *ElPersona) createfunctionid() error {
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

func (s *ElPersona) dropfunctionid() error {
	query := `
    drop function if exists fn_trig_persona_pk;
    `
	err := s.db.Exec(query)
	return err
}

func (s *ElPersona) createtriggerid() error {
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

func (s *ElPersona) droptrigid() error {
	query := `
    drop trigger if exists trig_persona_pk on persona;
    `
	err := s.db.Exec(query)
	return err
}

func (s *ElPersona) InsertPersona(elpersona *Persona) error {
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

func (s *ElPersona) DeletePersona(int) error {
	return nil
}

func (s *ElPersona) UpdatePersona(*Persona) error {
	return nil
}

// func (s *ElPersona) GetPersonaById(int) (*Persona, error) {
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

func (s *ElPersona) GetPersonasByUserId(id int) ([]*PersonaView, error) {
	var personas []*PersonaView

	err := s.db.QueryScan(&personas, `select id, name from Persona where userid = $1`, id)
	if err != nil {
		return nil, err
	}
	if len(personas) == 0 {
		return nil, s.db.NotFound
	}

	return personas, nil
}
