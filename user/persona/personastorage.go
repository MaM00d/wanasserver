package persona

import (
	"database/sql"
	"fmt"
	"log/slog"

	_ "github.com/lib/pq"
)

func (s *ElPersona) InitDb() error {
	if err := s.createPersonaTabel(); err != nil {
		return err
	}
	if err := s.createuserfk(); err != nil {
		return err
	}

	if err := s.createfunctionid(); err != nil {
		return err
	}
	if err := s.createtriggerid(); err != nil {
		return err
	}

	return s.createtriggerid()
}

func (s *ElPersona) DropDb() error {
	if err := s.droptrigid(); err != nil {
		return err
	}
	if err := s.dropfunctionid(); err != nil {
		return err
	}
	if err := s.dropPersonaTabel(); err != nil {
		return err
	}
	if err := s.dropuserfk(); err != nil {
		return err
	}

	return nil
}

func (s *ElPersona) dropPersonaTabel() error {
	query := `
    drop table if exists Persona;
    `
	err := s.db.Exec(query)
	return err
}

func (s *ElPersona) dropuserfk() error {
	query := `
    drop CONSTRAINT if exists fk_persona_userid
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

func (s *ElPersona) dropfunctionid() error {
	query := `
    drop function if exists fn_trig_persona_pk;
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

func (s *ElPersona) createPersonaTabel() error {
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

func (s *ElPersona) GetPersonasByUserId(id int) ([]Persona, error) {
	var personas []Persona
	rows := s.db.QueryScan(personas, `select * from Personas where ID = $1`, id)

	if rows == fmt.Errorf("not found") {

		slog.Error("GetPersonasByUserId", "id", id)
		return nil, fmt.Errorf("persona with id [%d] not found", id)
	}
	return personas, nil
}
