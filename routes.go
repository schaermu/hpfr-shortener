package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type URLShortenRequest struct {
	URL string `json:"url"`
}

type URLShortenResponse struct {
	ShortURL string `json:"short_url"`
}

func Shorten(c echo.Context) (err error) {
	url := new(URLShortenRequest)
	if err = c.Bind(url); err != nil {
		return
	}

	return c.JSON(http.StatusOK, url)
}
