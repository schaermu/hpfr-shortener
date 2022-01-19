package main

import (
	"context"

	"github.com/schaermu/hpfr-shortener/cmd"
	"github.com/sirupsen/logrus"

	_ "github.com/joho/godotenv/autoload"
)

var log = logrus.New()

func main() {
	server, err := cmd.NewServer(getFileSystem(), log)
	if err != nil {
		log.Fatal(err)
	}

	// make sure the mongodb connection is closed on shutdown
	defer func() {
		if err := server.Datastore.Session.Disconnect(context.Background()); err != nil {
			panic(err)
		}
	}()

	log.Fatal(server.Http.Start(":8080"))
}
