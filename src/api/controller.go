package api

import "github.com/labstack/echo"

type Controller interface {
	RegisterRoutes(group echo.Group)
}
