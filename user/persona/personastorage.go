package persona

import (
	"database/sql"
	"log/slog"

	_ "github.com/lib/pq"
)

func (s *ElPersona) createPersonaTabel() error {
	query := `
          CREATE TABLE IF NOT EXISTS persona (id int NOT NULL, name char(50) NOT NULL, userid int NOT NULL, createdat timestamp NOT NULL, CONSTRAINT pk_persona PRIMARY KEY (id
                                                                                                                                                                               , userid));
    `
	err := s.db.Exec(query)
	return err
}

func (s *ElPersona) dropPersonaTabel() error {
	query := `
          DROP TABLE IF EXISTS persona;
    `
	err := s.db.Exec(query)
	return err
}

func (s *ElPersona) createuserfk() error {
	query := `
          ALTER TABLE persona ADD CONSTRAINT fk_persona_userid
            FOREIGN key(userid) REFERENCES users (id);
    `
	err := s.db.Exec(query)
	return err
}

func (s *ElPersona) dropuserfk() error {
	query := `
          ALTER TABLE IF EXISTS persona
            DROP CONSTRAINT IF EXISTS fk_persona_userid;
    `
	err := s.db.Exec(query)
	return err
}

func (s *ElPersona) createfunctionid() error {
	query := `
          CREATE OR REPLACE FUNCTION "fn_trig_persona_pk"() RETURNS "pg_catalog"."trigger" AS $BODY$
                      begin
                      new.id = (select count(*) from persona where userid=new.userid);
                      return NEW;
                      end;
                      $BODY$ LANGUAGE PLPGSQL VOLATILE cost 100;
    `
	err := s.db.Exec(query)
	return err
}

func (s *ElPersona) dropfunctionid() error {
	query := `
          DROP FUNCTION IF EXISTS fn_trig_persona_pk;
    `
	err := s.db.Exec(query)
	return err
}

func (s *ElPersona) createtriggerid() error {
	query := `
          CREATE TRIGGER trig_persona_pk
            BEFORE
            INSERT ON persona
            FOR EACH ROW EXECUTE PROCEDURE fn_trig_persona_pk();
    `
	err := s.db.Exec(query)
	return err
}

func (s *ElPersona) droptrigid() error {
	query := `
          DROP TRIGGER IF EXISTS trig_persona_pk ON persona;
    `
	err := s.db.Exec(query)
	return err
}

func (s *ElPersona) InsertPersona(elpersona *Persona) error {
	query := `insert into Persona 
          (name,userid,createdat)
            VALUES ($1,$2,$3)
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
