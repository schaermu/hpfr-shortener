package utils

import (
	"fmt"

	"github.com/labstack/echo/v4"
)

func GetShortLink(id string, ctx echo.Context, cfg Config) string {
	if len(cfg.BaseURL) > 0 {
		return fmt.Sprintf("%s/%s", cfg.BaseURL, id)
	}
	return fmt.Sprintf("%s://%s/%s", ctx.Scheme(), ctx.Request().Host, id)
}
