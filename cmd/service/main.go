package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/echo-jwt/v4"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"billiard_app_backend/handlers"
	"billiard_app_backend/models"
)

var db *gorm.DB

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	var signingKey = []byte("superdupersecret!")
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set(models.SigningContextKey, signingKey)
			return next(c)
		}
	})

	var err error
	dsn := "host=db user=postgres password=postgres dbname=notes port=5432 sslmode=disable"
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		e.Logger.Fatal(err)
	}

    if err := db.AutoMigrate(); err != nil {
        e.Logger.Fatal(err)
    }

	e.GET("/health-check", handlers.HealthCheck)

	v1 := e.Group("/v1")
	v1.POST("/login", handlers.Login)
	v1.POST("/logout", handlers.Logout)
	v1Notes := v1.Group("/notes", echojwt.JWT(signingKey))
	v1Notes.GET("", handlers.GetNote)
	v1Notes.PUT("", handlers.UpdateNote)

	e.Logger.Fatal(e.Start(":3000"))
}
