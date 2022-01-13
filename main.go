package main

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/schaermu/hpfr-shortener/internal/data"
	"github.com/schaermu/hpfr-shortener/internal/handlers"
	"github.com/schaermu/hpfr-shortener/internal/repositories"
	"github.com/schaermu/hpfr-shortener/internal/utils"
	"github.com/sirupsen/logrus"

	_ "github.com/joho/godotenv/autoload"
)

var log = logrus.New()

func main() {
	config := utils.NewConfigFromEnv()

	// connect datastore, setup repositories
	ds := data.NewDatastore(config.MongoDSN, config.MongoDB, log)
	defer func() {
		if err := ds.Session.Disconnect(context.Background()); err != nil {
			panic(err)
		}
	}()

	urlRepo := repositories.NewURLRepository(ds)

	e := echo.New()

	// setup middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Filesystem: getFileSystem(),
		HTML5:      true,
	}))

	// setup routes
	handlers.NewURLHandler(e, urlRepo)

	// start
	e.Logger.Fatal(e.Start(":8080"))
}
