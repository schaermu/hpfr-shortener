package repositories

import (
	"context"
	"time"

	"github.com/labstack/gommon/log"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/schaermu/hpfr-shortener/internal/data"
	"github.com/schaermu/hpfr-shortener/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type URLRepository struct {
	store *data.MongoDatastore
}

func NewURLRepository(store *data.MongoDatastore) *URLRepository {
	return &URLRepository{
		store: store,
	}
}

func (r *URLRepository) FindByShortCode(code string) (shortURL domain.ShortURL, err error) {
	err = r.store.URLCollection.FindOne(context.TODO(), bson.M{"code": code}).Decode(&shortURL)
	return
}

func (r *URLRepository) NewShortURL(url string) (string, error) {
	shortID, err := r.getUniqueID()
	if err != nil {
		return "", err
	}

	shortURL := domain.ShortURL{
		CreatedAt:     time.Now(),
		TargetURL:     url,
		ShortCode:     shortID,
		RedirectCount: 0,
	}
	res, err := r.store.URLCollection.InsertOne(context.TODO(), shortURL)
	if err != nil {
		return "", err
	}
	var objID = res.InsertedID.(primitive.ObjectID).Hex()
	log.Infof("Created new short url with id %v", objID)
	return objID, nil
}

func (r *URLRepository) getUniqueID() (string, error) {
	var found = ""
	for len(found) == 0 {
		id, err := gonanoid.New(6)
		if err != nil {
			return "", err
		}

		res, err := r.FindByShortCode(id)
		if res.ShortCode == "" {
			found = id
		}
	}
	return found, nil
}
