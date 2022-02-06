package cmd

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/schaermu/hpfr-shortener/internal/data"
	"github.com/schaermu/hpfr-shortener/internal/handlers"
	"github.com/schaermu/hpfr-shortener/internal/repositories"
	"github.com/schaermu/hpfr-shortener/internal/utils"
	"github.com/sirupsen/logrus"
)

type Server struct {
	Http      *echo.Echo
	Datastore *data.MongoDatastore
}

func NewServer(embedFS http.FileSystem, log *logrus.Logger) (server *Server, err error) {
	config, err := utils.NewConfigFromEnv()
	if err != nil {
		return nil, err
	}

	// connect datastore, setup repositories
	ds, err := data.NewDatastore(config.MongoDSN, config.MongoDB, log)
	if err != nil {
		return nil, err
	}
	urlRepo := repositories.NewURLRepository(ds, log)

	e := echo.New()

	// setup middleware
	e.Use(middleware.RequestID())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format:           "time=${time_custom}, id=${id}, method=${method}, uri=${uri}, status=${status},${error} lat=${latency_human}, bytes_out=${bytes_out}\n",
		CustomTimeFormat: "2006-01-02 15:04:05",
	}))
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOriginFunc: func(origin string) (bool, error) {
			if len(config.BaseURL) > 0 && origin == config.BaseURL {
				return true, nil
			}
			return origin == "http://localhost", nil
		},
		AllowMethods: []string{"GET", "POST"},
	}))
	e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Filesystem: embedFS,
		HTML5:      false,
	}))

	// register custom validator
	e.Validator = &utils.CustomValidator{Validator: validator.New()}

	// setup routes
	handlers.NewURLHandler(e, urlRepo, config)

	return &Server{
		Http:      e,
		Datastore: ds,
	}, nil
}
