package handlers

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	jwt "github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"billiard_app_backend/bindings"
	"billiard_app_backend/models"
	"billiard_app_backend/renderings"
)

func Login(c echo.Context) error {
	resp := renderings.LoginResponse{}
	req := new(bindings.LoginRequest)

	if err := c.Bind(req); err != nil {
		resp.Success = false
		resp.Message = "Unabled to bind request for login"

		return c.JSON(http.StatusBadRequest, resp)
	}

	if err := req.Validate(c); err != nil {
		resp.Success = false
		resp.Message = err.Error()

		return c.JSON(http.StatusBadRequest, resp)
	}

	db := c.Get(models.DBContextKey).(*gorm.DB)
	user, err := models.GetUserByUsername(db, req.Username)

	if err != nil {
		resp.Success = false
		resp.Message = "Username or Password incorrect"

		return c.JSON(http.StatusUnauthorized, resp)
	}

	if err := bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(req.Password)); err == nil {
		resp.Success = false
		resp.Message = "Username or Password incorrect"

		return c.JSON(http.StatusUnauthorized, resp)
	}

	signingKey := c.Get(models.SigningContextKey).([]byte)

	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		Issuer:    "service",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(signingKey)
	if err != nil {
		resp.Success = false
		resp.Message = "Server Error"

		return c.JSON(http.StatusInternalServerError, resp)
	}

	resp.Token = ss

	return c.JSON(http.StatusOK, resp)
}
