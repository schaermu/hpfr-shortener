package repositories

import (
	"context"
	"errors"
	"net/url"
	"time"

	ip2location "github.com/ip2location/ip2location-go"
	"github.com/labstack/echo/v4"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/schaermu/hpfr-shortener/internal/data"
	"github.com/schaermu/hpfr-shortener/internal/domain"
	"github.com/schaermu/hpfr-shortener/internal/utils"
	"github.com/sirupsen/logrus"
	"github.com/ua-parser/uap-go/uaparser"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type URLRepository struct {
	store    *data.MongoDatastore
	ip2locDB *ip2location.DB
	logger   *logrus.Logger
}

type FindOptions struct {
	IncludeHits bool
}

type TotalCount struct {
	Id    primitive.ObjectID `bson:"_id,omitempty"`
	Count int64              `bson:"count,omitempty"`
}

type TimeBasedStatistic struct {
	Date  time.Time `bson:"_id,omitempty"`
	Value int64     `bson:"value,omitempty"`
}

type StatisticsResult struct {
	TotalCount int64
	Hits       []TimeBasedStatistic
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
	return r.FindByIDWithOptions(id, &FindOptions{})
}

func (r *URLRepository) FindByIDWithOptions(id string, opts *FindOptions) (shortURL domain.ShortURL, err error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return
	}
	err = r.store.URLCollection.FindOne(context.TODO(),
		bson.M{"_id": bson.M{"$eq": objID}},
		options.FindOne().SetProjection(bson.M{"hits": utils.BTOI(opts.IncludeHits)})).Decode(&shortURL)
	return
}

func (r *URLRepository) FindByShortCode(code string) (shortURL domain.ShortURL, err error) {
	return r.FindByShortCodeWithOptions(code, &FindOptions{})
}

func (r *URLRepository) FindByShortCodeWithOptions(code string, opts *FindOptions) (shortURL domain.ShortURL, err error) {
	err = r.store.URLCollection.FindOne(context.TODO(),
		bson.M{"short_code": code},
		options.FindOne().SetProjection(bson.M{"hits": utils.BTOI(opts.IncludeHits)})).Decode(&shortURL)
	return
}

func (r *URLRepository) getHitsTimeSeries(code string) (result []TimeBasedStatistic, err error) {
	matchStage := bson.M{"$match": bson.M{"short_code": code}}
	unwindStage := bson.M{"$unwind": "$hits"}
	addFieldsStage := bson.M{"$addFields": bson.M{
		"hitDate": bson.M{
			"$dateFromParts": bson.M{
				"year":  bson.M{"$year": "$hits.created_at"},
				"month": bson.M{"$month": "$hits.created_at"},
				"day":   bson.M{"$dayOfMonth": "$hits.created_at"},
			},
		},
	}}
	groupStage := bson.M{"$group": bson.D{
		{Key: "_id", Value: "$hitDate"},
		{Key: "value", Value: bson.M{"$sum": 1}},
	}}

	hitsByDate, err := r.store.URLCollection.Aggregate(context.TODO(), []bson.M{
		unwindStage, matchStage, addFieldsStage, groupStage,
	})
	if err != nil {
		return
	}

	var hits []TimeBasedStatistic
	if err = hitsByDate.All(context.TODO(), &hits); err == nil {
		result = hits
	}
	return
}

func (r *URLRepository) getHitCount(code string) (result int64, err error) {
	matchStage := bson.M{"$match": bson.M{"short_code": code}}
	projectStage := bson.M{"$project": bson.M{"count": bson.M{"$size": "$hits"}}}

	hitCount, err := r.store.URLCollection.Aggregate(context.TODO(), []bson.M{matchStage, projectStage})
	if err != nil {
		return
	}

	hitCount.Next(context.TODO())
	var res = TotalCount{}
	if err = hitCount.Decode(&res); err != nil {
		return
	}
	result = res.Count
	return
}

func (r *URLRepository) GetStatistics(code string) (result StatisticsResult, err error) {
	result = StatisticsResult{}
	if timeSeriesHits, err := r.getHitsTimeSeries(code); err == nil {
		result.Hits = timeSeriesHits
	}

	if totalCount, err := r.getHitCount(code); err == nil {
		result.TotalCount = totalCount
	}

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

	r.logger.Infof("Created new short url with code %v", shortURL.ShortCode)
	return shortURL.ShortCode, nil
}

func (r *URLRepository) RecordHit(target domain.ShortURL, c echo.Context) error {

	// parse user agent
	parser := uaparser.NewFromSaved()
	uap := parser.Parse(c.Request().UserAgent())

	hit := domain.ShortURLHit{
		CreatedAt:    time.Now(),
		UAFamily:     uap.UserAgent.Family,
		UAMajor:      uap.UserAgent.Major,
		UAMinor:      uap.UserAgent.Minor,
		OS:           uap.Os.Family,
		OSMajor:      uap.Os.Major,
		OSMinor:      uap.Os.Minor,
		DeviceFamily: uap.Device.Family,
		DeviceModel:  uap.Device.Model,
		DeviceBrand:  uap.Device.Brand,
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
