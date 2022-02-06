package repositories

import (
	"context"
	"errors"
	"net/url"
	"time"

	ip2location "github.com/ip2location/ip2location-go"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/schaermu/hpfr-shortener/internal/data"
	"github.com/schaermu/hpfr-shortener/internal/domain"
	"github.com/sirupsen/logrus"
	"github.com/ua-parser/uap-go/uaparser"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type URLRepository struct {
	store    *data.MongoDatastore
	ip2locDB *ip2location.DB
	logger   *logrus.Logger
}

func NewURLRepository(store *data.MongoDatastore, logger *logrus.Logger) *URLRepository {
	var dbPath = "./IP2LOCATION-LITE-DB1.BIN"
	db, err := ip2location.OpenDB(dbPath)
	if err != nil {
		logger.Infof("Could not locate IP2Location database at %q, skipping ip lookups.", dbPath)
	}

	return &URLRepository{
		store:    store,
		ip2locDB: db,
		logger:   logger,
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
		CreatedAt: time.Now(),
		TargetURL: uri,
		ShortCode: shortID,
		Hits:      []domain.ShortURLHit{},
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

func (r *URLRepository) RecordHit(target domain.ShortURL, c echo.Context) error {

	// parse user agent
	parser := uaparser.NewFromSaved()
	uap := parser.Parse(c.Request().UserAgent())

	hit := domain.ShortURLHit{
		CreatedAt: time.Now(),
		UAFamily:  uap.UserAgent.Family,
		UAMajor:   uap.UserAgent.Major,
		UAMinor:   uap.UserAgent.Minor,
		OS:        uap.Os.Family,
		OSMajor:   uap.Os.Major,
		OSMinor:   uap.Os.Minor,
	}

	// try to locate ip address
	if r.ip2locDB != nil {
		if res, err := r.ip2locDB.Get_all(c.RealIP()); err == nil {
			hit.Country = res.Country_long
			hit.CountryCode = res.Country_short
		}
	}

	_, err := r.store.URLCollection.UpdateByID(context.TODO(), target.ID, bson.M{
		"$push": bson.M{"hits": hit},
	})

	if err != nil {
		return err
	}

	return nil
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
