package main

import (
	"context"
	"embed"
	"io/fs"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/schaermu/hpfr-shortener/internal/data"
	"github.com/schaermu/hpfr-shortener/internal/handlers"
	"github.com/schaermu/hpfr-shortener/internal/repositories"
	"github.com/schaermu/hpfr-shortener/internal/utils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

//go:embed ui/dist
var embeddedFiles embed.FS

var log = logrus.New()

func getFileSystem() http.FileSystem {
	fsys, err := fs.Sub(embeddedFiles, "ui/dist")
	if err != nil {
		panic(err)
	}

	return http.FS(fsys)
}

func main() {
	// read config
	config, err := loadConfig(".")
	if err != nil {
		log.Fatal("cannot load config: ", err)
	}

	// connect datastore, setup repositories
	ds := data.NewDatastore(config.MongoDSN, config.MongoDB, log)
	defer func() {
		if err = ds.Session.Disconnect(context.Background()); err != nil {
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

func loadConfig(path string) (config utils.Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
