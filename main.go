package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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

	db.AutoMigrate(&Note{})

	e.GET("/notes", getNotes)
	e.POST("/notes/new", createNote)
	e.GET("/notes/:id", getNote)
	e.PUT("/notes/:id", updateNote)
	e.DELETE("/notes/:id", deleteNote)

	e.Logger.Fatal(e.Start(":3000"))
}

type Note struct {
	ID      uint   `json:"id" gorm:"primaryKey"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

func getNotes(c echo.Context) error {
	var notes []Note
	if err := db.Find(&notes).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Error retrieving notes"})
	}
	return c.JSON(http.StatusOK, notes)
}

func createNote(c echo.Context) error {
	note := new(Note)
	if err := c.Bind(note); err != nil {
		return err
	}
	db.Create(note)
	return c.JSON(http.StatusCreated, note)
}

func getNote(c echo.Context) error {
	id := c.Param("id")
	var note Note
	if err := db.First(&note, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Note not found"})
	}
	return c.JSON(http.StatusOK, note)
}

func updateNote(c echo.Context) error {
	id := c.Param("id")
	var note Note
	if err := db.First(&note, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Note not found"})
	}
	if err := c.Bind(&note); err != nil {
		return err
	}
	db.Save(&note)
	return c.JSON(http.StatusOK, note)
}

func deleteNote(c echo.Context) error {
	id := c.Param("id")
	if err := db.Delete(&Note{}, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Note not found"})
	}
	return c.NoContent(http.StatusNoContent)
}
