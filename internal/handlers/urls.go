package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/schaermu/hpfr-shortener/internal/repositories"
	"github.com/schaermu/hpfr-shortener/internal/utils"
)

type URLShortenRequest struct {
	URL string `json:"url" validate:"required,url"`
}

type URLShortenResponse struct {
	ShortURL string `json:"short_url"`
}

type URLHandler struct {
	repository *repositories.URLRepository
	config     *utils.Config
}

func NewURLHandler(e *echo.Echo, repository *repositories.URLRepository, config *utils.Config) *URLHandler {
	handler := &URLHandler{
		repository: repository,
		config:     config,
	}

	e.POST("/api/shorten", handler.Shorten)
	e.GET("/:code", handler.Redirect)
	return handler
}

func (h *URLHandler) Shorten(c echo.Context) (err error) {
	url := new(URLShortenRequest)
	if err = c.Bind(url); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err = c.Validate(url); err != nil {
		return err
	}

	id, err := h.repository.NewShortURL(url.URL)
	if err != nil {
		return
	}

	var shortURL = utils.GetShortLink(id, c, *h.config)
	return c.JSON(http.StatusCreated, &URLShortenResponse{ShortURL: shortURL})
}

func (h *URLHandler) Redirect(c echo.Context) (err error) {
	code := c.Param("code")
	target, err := h.repository.FindByShortCode(code)
	if err != nil {
		return echo.NewHTTPError(404)
	}
	return c.Redirect(302, target.TargetURL)
}
