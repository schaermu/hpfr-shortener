package data

import (
	"context"
	"time"

	"github.com/pkg/errors"
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

func NewDatastore(dsn string, dbname string, logger *logrus.Logger) (ds *MongoDatastore, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dsn))

	if err == nil {
		db := client.Database(dbname)
		ds = new(MongoDatastore)
		ds.URLCollection = db.Collection(urlCollectionName)
		ds.logger = logger
		ds.Session = client

		err = client.Ping(context.TODO(), nil)
		if err != nil {
			return nil, err
		}

		logger.Infof("Connected to database %v/%v", dsn, dbname)
		return
	}

	return nil, errors.Errorf("error connecting to database: %v/%v", dsn, dbname)
}
