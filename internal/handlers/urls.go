package handlers

import (
	"net/http"
	"strings"

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

type URLStatisticsResponse struct {
	HitCount    int64     `json:"hits"`
	HitTimeData [][]int64 `json:"hitTimeData`
}

type URLHandler struct {
	repository *repositories.URLRepository
	config     *utils.Config
}

var StaticFS http.FileSystem

func NewURLHandler(e *echo.Echo, repository *repositories.URLRepository, config *utils.Config) *URLHandler {
	handler := &URLHandler{
		repository: repository,
		config:     config,
	}

	e.POST("/api/shorten", handler.Shorten)
	e.GET("/api/statistics", handler.Statistics)
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

	target, err := h.repository.FindByShortCode(strings.Trim(code, "+"))
	if err != nil {
		c.Logger().Debugf("%v", err)
		return echo.NewHTTPError(http.StatusNotFound)
	}

	// requests to statistics are handled here
	if strings.HasSuffix(code, "+") {
		return h.RenderStatistics(c)
	}

	defer h.repository.RecordHit(target, c)

	return c.Redirect(http.StatusTemporaryRedirect, target.TargetURL)
}

func (h *URLHandler) RenderStatistics(c echo.Context) (err error) {
	f, err := StaticFS.Open("index.html")
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error serving stats")
	}

	return c.Stream(http.StatusFound, "text/html", f)
}

func (h *URLHandler) Statistics(c echo.Context) (err error) {
	code := c.QueryParam("code")
	if len(code) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	// create statistics
	stats, err := h.repository.GetStatistics(code)
	if err != nil {
		return echo.NewHTTPError(500, err)
	}

	var hitTimeData = [][]int64{}
	for i := 0; i < len(stats.Hits); i++ {
		hitTimeData = append(hitTimeData, []int64{
			stats.Hits[i].Date.UnixMilli(),
			stats.Hits[i].Value,
		})
	}

	return c.JSON(http.StatusFound, &URLStatisticsResponse{
		HitTimeData: hitTimeData,
	})
}
