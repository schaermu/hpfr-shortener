package handlers

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/schaermu/hpfr-shortener/internal/repositories"
)

type URLShortenRequest struct {
	URL string `json:"url"`
}

type URLShortenResponse struct {
	ShortURL string `json:"short_url"`
}

type URLHandler struct {
	repository *repositories.URLRepository
}

func NewURLHandler(e *echo.Echo, repository *repositories.URLRepository) {
	handler := &URLHandler{
		repository: repository,
	}

	e.POST("/api/shorten", handler.Shorten)
}

func (h *URLHandler) Shorten(c echo.Context) (err error) {
	url := new(URLShortenRequest)
	if err = c.Bind(url); err != nil {
		return
	}

	id, err := h.repository.NewShortURL(url.URL)
	if err != nil {
		return
	}

	return c.JSON(http.StatusOK, &URLShortenResponse{ShortURL: fmt.Sprintf("https://hpfr.ch/%s", id)})
}
