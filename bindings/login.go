package bindings

import "github.com/labstack/echo/v4"

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (req LoginRequest) Validate(c echo.Context) error {
	errs := new(RequestErrors)
	if req.Username == "" {
		errs.Append(ErrUsernameEmpty)
	}
	if req.Password == "" {
		errs.Append(ErrPasswordEmpty)
	}
	if errs.Len() == 0 {
		return nil
	}
	return errs
}
