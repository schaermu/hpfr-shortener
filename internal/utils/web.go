package utils

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

func GetShortLink(id string, ctx echo.Context, cfg Config) string {
	if len(cfg.BaseURL) > 0 {
		return fmt.Sprintf("%s/%s", cfg.BaseURL, id)
	}
	return fmt.Sprintf("%s://%s/%s", ctx.Scheme(), ctx.Request().Host, id)
}

type CustomValidator struct {
	Validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.Validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}
