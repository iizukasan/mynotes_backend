package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"billiard_app_backend/renderings"
)

func HealthCheck(c echo.Context) error {
	resp := renderings.HealthCheckResponse{
		Message: "Everything is good!",
	}
	return c.JSON(http.StatusOK, resp)
}
