package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"billiard_app_backend/handlers"
)

var db *gorm.DB

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

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

	e.POST("/login", handlers.Login)
	e.POST("/logout", handlers.Logout)

	e.Logger.Fatal(e.Start(":3000"))
}
