package main

import (
	"hochzeit/src/api"
	"hochzeit/src/config"
	"hochzeit/src/pg/migrations"
	"hochzeit/src/service"

	"github.com/labstack/echo/middleware"

	"github.com/jmoiron/sqlx"

	"github.com/labstack/echo"
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())

	// cfg & db
	cfg := config.NewConfig()
	db := sqlx.MustConnect("postgres", cfg.PostgresConfig.GetConnectionURL())
	migrations.MustMigrate(db)

	// app context
	validator := api.NewValidator()
	rsvpService := service.NewRsvpService(db)
	rsvpController := api.NewRsvpController(validator, rsvpService)
	adminController := api.NewAdminController(rsvpService)

	// static routes
	e.Static("/static", "./static")
	e.Static("/", "./static/index.html")
	e.Use(CacheControl)

	// api routes
	api := e.Group("/api")
	rsvpController.RegisterRoutes(api)

	// admin
	admin := e.Group("/admin")
	admin.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (authenticated bool, err error) {
		authenticated = (username == cfg.BasicAuthConfig.Username && password == cfg.BasicAuthConfig.Password)
		return
	}))

	admin.Static("", "./static/admin.html")
	adminController.RegisterRoutes(admin)

	if err := e.Start(":8080"); err != nil {
		panic(err)
	}
}

func CacheControl(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set("Cache-Control", "no-cache")
		return next(c)
	}
}
