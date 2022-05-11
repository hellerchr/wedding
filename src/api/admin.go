package api

import (
	"hochzeit/src/service"

	"github.com/labstack/echo"
	_ "github.com/lib/pq" //pg driver
)

type AdminController struct {
	service *service.RsvpService
}

func NewAdminController(service *service.RsvpService) *AdminController {
	return &AdminController{service: service}
}

func (r *AdminController) RegisterRoutes(group *echo.Group) {
	group.DELETE("/rsvp/:id", r.delete)
	group.GET("/rsvp", r.list)
}

func (r *AdminController) delete(c echo.Context) (err error) {
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(400, "no id provided")
	}

	err = r.service.Delete(id)
	if err != nil {
		return err
	}

	return c.String(200, "")
}

func (r *AdminController) list(c echo.Context) (err error) {
	rsvps, err := r.service.List()
	if err != nil {
		return err
	}

	return c.JSON(200, rsvps)
}
