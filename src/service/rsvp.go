package service

import (
	"hochzeit/src/pg"
	"time"

	"github.com/labstack/echo"

	"github.com/jmoiron/sqlx"
)

type Rsvp struct { //nolint: maligned
	ID          string `db:"id"`
	SessionID   string `db:"session_id"`
	Name        string `db:"name"`
	Email       string `db:"email"`
	Accept      bool   `db:"accept"`
	PersonCount int    `db:"person_count"`
	PersonNames []string
	Diet        string    `db:"diet"`
	Hotel       bool      `db:"hotel"`
	Message     string    `db:"message"`
	CreatedOn   time.Time `db:"created_on"`
}

type RsvpName struct { //nolint: maligned
	RsvpID string `db:"rsvp_id"`
	Name   string `db:"name"`
}

type RsvpService struct {
	db *sqlx.DB
}

func NewRsvpService(db *sqlx.DB) *RsvpService {
	return &RsvpService{db: db}
}

func (r *RsvpService) Delete(id string) (err error) {
	sb := pg.NewStatementBuilder(r.db)
	result, err := sb.Delete("rsvp").Where("id = ?", id).Exec()

	if err != nil {
		return err
	}

	i, _ := result.RowsAffected()
	if i == 0 {
		return echo.NewHTTPError(404, "entry not found")
	}

	_, err = sb.Delete("rsvp_names").Where("rsvp_id = ?", id).Exec()

	return
}

func (r *RsvpService) List() (rsvps []Rsvp, err error) {
	rows, err := r.db.Queryx("SELECT * FROM rsvp")
	for rows.Next() {
		var r Rsvp
		if err = rows.StructScan(&r); err != nil {
			return
		}
		rsvps = append(rsvps, r)
	}

	return
}

func (r *RsvpService) Create(rsvp Rsvp) (err error) {
	_, err = pg.Insert(r.db, "rsvp", &rsvp)

	for _, name := range rsvp.PersonNames {
		_, err = pg.Insert(r.db, "rsvp_names", &RsvpName{
			RsvpID: rsvp.ID,
			Name:   name,
		})
		if err != nil {
			return
		}
	}

	return
}
