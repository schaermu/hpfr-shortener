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

func (r *URLRepository) GetAll() (shortURLs []domain.ShortURL, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cur, err := r.store.URLCollection.Find(ctx, bson.D{})
	defer cur.Close(ctx)
	if err != nil {
		return nil, err
	}

	var res = []domain.ShortURL{}
	for cur.Next(ctx) {
		var result domain.ShortURL
		if err := cur.Decode(&result); err != nil {
			log.Fatal(err)
		}
		res = append(res, result)
	}
	return res, nil
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
