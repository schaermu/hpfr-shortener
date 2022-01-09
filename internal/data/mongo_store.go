package data

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	urlCollectionName = "urls"
)

type MongoDatastore struct {
	URLCollection *mongo.Collection
	Session       *mongo.Client
	logger        *logrus.Logger
}

func NewDatastore(dsn string, dbname string, logger *logrus.Logger) *MongoDatastore {
	var ds *MongoDatastore

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dsn))

	if err == nil {
		db := client.Database(dbname)
		ds = new(MongoDatastore)
		ds.URLCollection = db.Collection(urlCollectionName)
		ds.logger = logger
		ds.Session = client

		logger.Infof("Connected to database %v/%v", dsn, dbname)

		return ds
	}

	logger.Fatalf("Error connecting to database: %v/%v", dsn, dbname)
	return nil
}
