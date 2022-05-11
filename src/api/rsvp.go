package api

import (
	"hochzeit/src/service"
	"net/http"
	"reflect"
	"strconv"
	"time"

	"github.com/labstack/echo"
	_ "github.com/lib/pq" //pg driver
	uuid "github.com/satori/go.uuid"
)

const namePerson = "name_person"

type RsvpRequest struct { //nolint: maligned
	Name        string `form:"name" validate:"required,min=3,max=50"`
	Email       string `form:"email" validate:"required,email"`
	Accept      bool   `form:"accept"`
	PersonCount int    `form:"person_count" validate:"min=0,max=8"`
	NamePerson2 string `form:"name_person2"`
	NamePerson3 string `form:"name_person3"`
	NamePerson4 string `form:"name_person4"`
	NamePerson5 string `form:"name_person5"`
	NamePerson6 string `form:"name_person6"`
	NamePerson7 string `form:"name_person7"`
	NamePerson8 string `form:"name_person8"`
	Diet        string `form:"diet" validate:"max=1000"`
	Hotel       bool   `form:"hotel"`
	Message     string `form:"message" validate:"max=2000"`
}

type RsvpController struct {
	validator *Validator
	service   *service.RsvpService
}

func (r *RsvpController) RegisterRoutes(group *echo.Group) {
	group.POST("/rsvp", r.create)
}

func NewRsvpController(validator *Validator, service *service.RsvpService) *RsvpController {
	return &RsvpController{validator: validator, service: service}
}

func (r *RsvpController) create(c echo.Context) (err error) {
	rsvp := new(RsvpRequest)
	if err = c.Bind(rsvp); err != nil {
		return err
	}

	if valErrs := validateCreate(r, rsvp); valErrs != nil {
		return c.JSON(400, valErrs)
	}

	sessionID := getRsvpSessionID(c)

	err = r.service.Create(service.Rsvp{
		ID:          uuid.NewV4().String(),
		SessionID:   sessionID,
		Name:        rsvp.Name,
		Email:       rsvp.Email,
		Accept:      rsvp.Accept,
		PersonCount: rsvp.PersonCount,
		PersonNames: getPersonNames(rsvp),
		Diet:        rsvp.Diet,
		Hotel:       rsvp.Hotel,
		Message:     rsvp.Message,
		CreatedOn:   time.Now(),
	})

	if err != nil {
		return err
	}

	return c.String(200, "")
}

func getPersonNames(rsvp *RsvpRequest) []string {
	all := []string{
		rsvp.Name,
		rsvp.NamePerson2,
		rsvp.NamePerson3,
		rsvp.NamePerson4,
		rsvp.NamePerson5,
		rsvp.NamePerson6,
		rsvp.NamePerson7,
		rsvp.NamePerson8,
	}

	return all[0:rsvp.PersonCount]
}

func getRsvpSessionID(c echo.Context) string { //nolint: unused, deadcode
	var rsvpSessionID string

	cookie, err := c.Cookie("rsvp_session_id")
	if err == nil && cookie.Value != "" {
		rsvpSessionID = cookie.Value
	} else {
		rsvpSessionID = uuid.NewV4().String()
		c.SetCookie(&http.Cookie{
			Name:  "rsvp_session_id",
			Value: rsvpSessionID,
		})
	}

	return rsvpSessionID
}

func validateCreate(r *RsvpController, rsvp *RsvpRequest) []ValidationError {
	validationErrors := r.validator.Validate(rsvp)

	// person count needs to be set when rsvp is accepted
	if rsvp.Accept && rsvp.PersonCount == 0 {
		validationErrors = append(validationErrors, ValidationError{
			Field:     "person_count",
			Violation: "custom",
		})
	}
	// validate person names
	for i := 2; i <= rsvp.PersonCount; i++ {
		elem := reflect.ValueOf(rsvp).Elem()
		formField := namePerson + strconv.Itoa(i)

		for i := 0; i < elem.NumField(); i++ {
			field := elem.Type().Field(i).Tag.Get("form")
			if field == formField {
				if len(elem.Field(i).String()) < 3 {
					validationErrors = append(validationErrors, ValidationError{
						Field:     formField,
						Violation: "custom",
					})
				}
			}
		}
	}

	return validationErrors
}
