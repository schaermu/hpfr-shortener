package repositories

import (
	"context"
	"errors"
	"net/url"
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

func (r *URLRepository) FindByID(id string) (shortURL domain.ShortURL, err error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return
	}
	err = r.store.URLCollection.FindOne(context.TODO(), bson.M{"_id": bson.M{"$eq": objID}}).Decode(&shortURL)
	return
}

func (r *URLRepository) FindByShortCode(code string) (shortURL domain.ShortURL, err error) {
	err = r.store.URLCollection.FindOne(context.TODO(), bson.M{"short_code": code}).Decode(&shortURL)
	return
}

func (r *URLRepository) NewShortURL(uri string) (string, error) {
	if uri == "" {
		return "", errors.New("url cannot be empty")
	}

	if u, err := url.Parse(uri); err != nil || (u.Scheme == "" || u.Host == "") {
		return "", errors.New("url is invalid")
	}

	shortID, err := r.getUniqueID()
	if err != nil {
		return "", err
	}

	shortURL := domain.ShortURL{
		CreatedAt:     time.Now(),
		TargetURL:     uri,
		ShortCode:     shortID,
		RedirectCount: 0,
	}

	res, err := r.store.URLCollection.InsertOne(context.TODO(), shortURL)
	if err != nil {
		return "", err
	}

	shortURL, err = r.FindByID(res.InsertedID.(primitive.ObjectID).Hex())
	if err != nil {
		return "", err
	}

	log.Infof("Created new short url with code %v", shortURL.ShortCode)
	return shortURL.ShortCode, nil
}

func (r *URLRepository) getUniqueID() (string, error) {
	var found = ""
	for len(found) == 0 {
		id, err := gonanoid.New(6)
		if err != nil {
			return "", err
		}

		res, _ := r.FindByShortCode(id)
		if res.ShortCode == "" {
			found = id
		}
	}
	return found, nil
}
