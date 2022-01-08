package main

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

//go:embed ui/dist
var embeddedFiles embed.FS

func getFileSystem() http.FileSystem {
	fsys, err := fs.Sub(embeddedFiles, "ui/dist")
	if err != nil {
		panic(err)
	}

	return http.FS(fsys)
}

func main() {
	e := echo.New()

	// middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Filesystem: getFileSystem(),
		HTML5:      true,
	}))

	// routes
	e.POST("/api/shorten", func(c echo.Context) (err error) {
		url := new(URLShortenRequest)
		if err = c.Bind(url); err != nil {
			return
		}

		return c.JSON(http.StatusOK, url)
	})

	// start
	e.Logger.Fatal(e.Start(":8080"))
}

type URLShortenRequest struct {
	URL string `json:"url"`
}
