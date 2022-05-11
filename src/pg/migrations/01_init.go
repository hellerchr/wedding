package migrations

import migrate "github.com/rubenv/sql-migrate"

func Migration01() *migrate.Migration {
	up := `
CREATE TABLE rsvp
(
	id varchar(36),
	session_id varchar(36),
    name     varchar(255),
    email    varchar(255),
    accept   boolean,
    person_count   smallint,
    diet     varchar(1000),
    hotel    boolean,
    message  varchar(2000),
	created_on timestamp NOT NULL
);
CREATE TABLE rsvp_names
(
	rsvp_id varchar(36),
    name     varchar(255)
)`

	down := "DROP TABLE rsvp, rsvp_names"

	return &migrate.Migration{
		Id:   "1",
		Up:   []string{up},
		Down: []string{down},
	}
}
