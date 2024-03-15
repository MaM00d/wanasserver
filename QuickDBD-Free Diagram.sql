-- Exported from QuickDBD: https://www.quickdatabasediagrams.com/
-- Link to schema: https://app.quickdatabasediagrams.com/#/d/rwe3bc
-- NOTE! If you have used non-SQL datatypes in your design, you will have to change these here.


CREATE TABLE users (
    id serial   NOT NULL,
    name char(50)   NOT NULL,
    email char(50)   NOT NULL,
    password char(100)   NOT NULL,
    phone int   NULL,
    createdAt timestamp   NOT NULL,
    CONSTRAINT pk_users PRIMARY KEY (
        id
     )
);

CREATE TABLE persona (
    id serial   NOT NULL,
    name char(50)   NOT NULL,
    useriD int   NOT NULL,
    createdat timestamp   NOT NULL,
    CONSTRAINT pk_persona PRIMARY KEY (
        id,useriD
     )
);

CREATE TABLE chat (
    id serial   NOT NULL,
    useriD int   NOT NULL,
    personaid int   NOT NULL,
    createdat timestamp   NOT NULL,
    CONSTRAINT pk_chat PRIMARY KEY (
        id,personaid,useriD
     )
);

CREATE TABLE msg (
    id serial   NOT NULL,
    useriD int   NOT NULL,
    chatid int   NOT NULL,
    personaid int   NOT NULL,
    message char(100)   NOT NULL,
    createdat timestamp   NOT NULL,
    CONSTRAINT pk_msg PRIMARY KEY (
        id,chatid
     )
);

ALTER TABLE persona ADD CONSTRAINT fk_persona_useriD FOREIGN KEY(useriD)
REFERENCES users (id);

ALTER TABLE chat ADD CONSTRAINT fk_chat_personaid FOREIGN KEY(personaid)
REFERENCES persona (id);

ALTER TABLE msg ADD CONSTRAINT fk_msg_chatid FOREIGN KEY(chatid)
REFERENCES chat (id);

CREATE INDEX idx_users_name
ON users (name);





CREATE OR REPLACE FUNCTION "fn_trig_pk"()
  RETURNS "persona"."trigger" AS $BODY$ 
begin
new.id = (select count(*)+1 from users where users.id=new.userid);
return NEW;
end;
