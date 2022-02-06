package repositories

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/schaermu/hpfr-shortener/internal/data"
	"github.com/schaermu/hpfr-shortener/internal/domain"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/tryvium-travels/memongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var logger = logrus.New()

type URLRepositoryTestSuite struct {
	suite.Suite
	mongoServer *memongo.Server
	store       *data.MongoDatastore
}

func (suite *URLRepositoryTestSuite) SetupSuite() {
	// load config from env
	mVersion := os.Getenv("MEMONGO_VERSION")
	if mVersion == "" {
		mVersion = "5.0.5"
	}

	// setup mongodb & repository
	mongoServer, err := memongo.Start(mVersion)
	if err != nil {
		suite.T().Fatal(err)
	}

	suite.store, err = data.NewDatastore(mongoServer.URI(), memongo.RandomDatabase(), logger)
	suite.mongoServer = mongoServer
}

func (suite *URLRepositoryTestSuite) TearDownSuite() {
	suite.mongoServer.Stop()
}

func (suite *URLRepositoryTestSuite) SetupTest() {
	suite.store.URLCollection.DeleteMany(context.TODO(), bson.M{})
}

func TestURLRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(URLRepositoryTestSuite))
}

func (suite *URLRepositoryTestSuite) TestFindByID() {
	for scenario, fn := range map[string]func(){
		"no id":    suite.testFindByID_NoID,
		"valid id": suite.testFindByID_ValidID,
	} {
		suite.SetupTest()
		suite.Run(scenario, fn)
	}
}

func (suite *URLRepositoryTestSuite) TestFindByShortCode() {
	for scenario, fn := range map[string]func(){
		"no code":    suite.testFindByShortCode_NoCode,
		"valid code": suite.testFindByShortCode_ValidCode,
	} {
		suite.SetupTest()
		suite.Run(scenario, fn)
	}
}

func (suite *URLRepositoryTestSuite) TestNewShortURL() {
	for scenario, fn := range map[string]func(){
		"no url":      suite.testNewShortURL_NoURL,
		"invalid url": suite.testNewShortURL_InvalidURL,
		"valid url":   suite.testNewShortURL_ValidURL,
	} {
		suite.SetupTest()
		suite.Run(scenario, fn)
	}
}

func (suite *URLRepositoryTestSuite) testFindByID_NoID() {
	// arrange
	var id = ""
	r := NewURLRepository(suite.store)

	// act
	_, err := r.FindByID(id)

	// assert
	assert.Error(suite.T(), err)
}

func (suite *URLRepositoryTestSuite) testFindByID_ValidID() {
	// arrange
	var (
		url  = "http://www.foobar.org"
		code = "foobar"
	)

	r := NewURLRepository(suite.store)
	var shortURL = domain.ShortURL{
		CreatedAt: time.Now(),
		TargetURL: url,
		ShortCode: code,
		Hits:      []domain.ShortURLHit{},
	}
	res, err := suite.store.URLCollection.InsertOne(context.TODO(), shortURL)
	var insertedID = res.InsertedID.(primitive.ObjectID).Hex()

	// act
	foundShortURL, err := r.FindByID(insertedID)

	// assert
	if assert.NoError(suite.T(), err) {
		assert.Equal(suite.T(), url, foundShortURL.TargetURL)
		assert.Equal(suite.T(), code, foundShortURL.ShortCode)
	}
}

func (suite *URLRepositoryTestSuite) testFindByShortCode_NoCode() {
	// arrange
	var code = ""
	r := NewURLRepository(suite.store)

	// act
	_, err := r.FindByShortCode(code)

	// assert
	assert.Error(suite.T(), err)
}

func (suite *URLRepositoryTestSuite) testFindByShortCode_ValidCode() {
	// arrange
	var (
		url  = "http://www.foobar.org"
		code = "foobar"
	)

	r := NewURLRepository(suite.store)
	var shortURL = domain.ShortURL{
		CreatedAt: time.Now(),
		TargetURL: url,
		ShortCode: code,
		Hits:      []domain.ShortURLHit{},
	}
	res, err := suite.store.URLCollection.InsertOne(context.TODO(), shortURL)
	var insertedID = res.InsertedID.(primitive.ObjectID).Hex()

	// act
	foundShortURL, err := r.FindByShortCode(code)

	// assert
	if assert.NoError(suite.T(), err) {
		assert.Equal(suite.T(), insertedID, foundShortURL.ID.Hex())
		assert.Equal(suite.T(), url, foundShortURL.TargetURL)
		assert.Equal(suite.T(), code, foundShortURL.ShortCode)
	}
}

func (suite *URLRepositoryTestSuite) testNewShortURL_NoURL() {
	// arrange
	var url = ""
	r := NewURLRepository(suite.store)

	// act
	_, err := r.NewShortURL(url)

	// assert
	assert.Error(suite.T(), err)
}

func (suite *URLRepositoryTestSuite) testNewShortURL_InvalidURL() {
	// arrange
	var url = "iNvAlId_CrAp"
	r := NewURLRepository(suite.store)

	// act
	_, err := r.NewShortURL(url)

	// assert
	assert.Error(suite.T(), err)
}

func (suite *URLRepositoryTestSuite) testNewShortURL_ValidURL() {
	// arrange
	var url = "http://www.foobar.org"
	r := NewURLRepository(suite.store)

	// act
	code, err := r.NewShortURL(url)

	// assert
	if assert.NoError(suite.T(), err) {
		shortURL, err := r.FindByShortCode(code)
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), code, shortURL.ShortCode)
		assert.Equal(suite.T(), url, shortURL.TargetURL)
	}
}
