package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/schaermu/hpfr-shortener/internal/data"
	"github.com/schaermu/hpfr-shortener/internal/repositories"
	"github.com/schaermu/hpfr-shortener/internal/utils"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/tryvium-travels/memongo"
)

var logger = logrus.New()

type URLHandlerTestSuite struct {
	suite.Suite
	echo        *echo.Echo
	mongoServer *memongo.Server
	repository  *repositories.URLRepository
}

func (suite *URLHandlerTestSuite) SetupSuite() {
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
	store := data.NewDatastore(mongoServer.URI(), memongo.RandomDatabase(), logger)
	suite.mongoServer = mongoServer
	suite.repository = repositories.NewURLRepository(store)

	// setup echo server
	suite.echo = echo.New()
	suite.echo.Validator = &utils.CustomValidator{Validator: validator.New()}
}

func (suite *URLHandlerTestSuite) TearDownSuite() {
	suite.mongoServer.Stop()
}

func TestURLHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(URLHandlerTestSuite))
}

func (suite *URLHandlerTestSuite) TestShorten() {
	for scenario, fn := range map[string]func(){
		"no url":      suite.testShortenEmptyURL,
		"invalid url": suite.testShortenInvalidURL,
		"proper url":  suite.testShortenProperURL,
	} {
		suite.Run(scenario, fn)
	}
}

func (suite *URLHandlerTestSuite) TestRedirect() {
	for scenario, fn := range map[string]func(){
		"empty id":        suite.testRedirectEmptyID,
		"non existent id": suite.testRedirectNonExistingID,
		"found":           suite.testRedirectFound,
	} {
		suite.Run(scenario, fn)
	}
}

func (suite *URLHandlerTestSuite) testRedirectEmptyID() {
	// arrange
	var id = ""
	rec := httptest.NewRecorder()
	c, h := suite.prepareRedirectTest(id, rec)

	// act
	err := h.Redirect(c)

	// assert
	if assert.Error(suite.T(), err) {
		if he, ok := err.(*echo.HTTPError); ok {
			assert.Equal(suite.T(), http.StatusNotFound, he.Code)
		}
	}
}

func (suite *URLHandlerTestSuite) testRedirectNonExistingID() {
	// arrange
	var id = "some_bogus_id"
	rec := httptest.NewRecorder()
	c, h := suite.prepareRedirectTest(id, rec)

	// act
	err := h.Redirect(c)

	if assert.Error(suite.T(), err) {
		if he, ok := err.(*echo.HTTPError); ok {
			assert.Equal(suite.T(), http.StatusNotFound, he.Code)
		}
	}
}

func (suite *URLHandlerTestSuite) testRedirectFound() {
	// arrange
	id, repoErr := suite.repository.NewShortURL("http://foobar.org")
	if repoErr != nil {
		suite.T().Fatal(repoErr)
	}
	rec := httptest.NewRecorder()
	c, h := suite.prepareRedirectTest(id, rec)

	// act
	err := h.Redirect(c)

	// assert
	if assert.NoError(suite.T(), err) {
		assert.Equal(suite.T(), http.StatusFound, rec.Code)
	}
}

func (suite *URLHandlerTestSuite) testShortenEmptyURL() {
	// arrange
	var url = ""
	rec := httptest.NewRecorder()
	c, h := suite.prepareShortenTest(url, rec)

	// act
	err := h.Shorten(c)

	// assert
	if assert.Error(suite.T(), err) {
		if he, ok := err.(*echo.HTTPError); ok {
			assert.Equal(suite.T(), http.StatusBadRequest, he.Code)
		}
	}
}

func (suite *URLHandlerTestSuite) testShortenInvalidURL() {
	// arrange
	var url = "bogus_url_input"
	rec := httptest.NewRecorder()
	c, h := suite.prepareShortenTest(url, rec)

	// act
	err := h.Shorten(c)

	// assert
	if assert.Error(suite.T(), err) {
		if he, ok := err.(*echo.HTTPError); ok {
			assert.Equal(suite.T(), http.StatusBadRequest, he.Code)
		}
	}
}

func (suite *URLHandlerTestSuite) testShortenProperURL() {
	// arrange
	var url = "https://foobar.org"
	rec := httptest.NewRecorder()
	c, h := suite.prepareShortenTest(url, rec)

	// act
	err := h.Shorten(c)

	// assert
	if assert.NoError(suite.T(), err) {
		assert.Equal(suite.T(), http.StatusCreated, rec.Code)
		var res = new(URLShortenResponse)
		err = json.Unmarshal([]byte(rec.Body.String()), &res)
		assert.NotEmpty(suite.T(), res.ShortURL)
	}
}

func (suite *URLHandlerTestSuite) prepareRedirectTest(id string, rec *httptest.ResponseRecorder) (c echo.Context, h *URLHandler) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	c = suite.echo.NewContext(req, rec)
	c.SetPath("/:code")
	c.SetParamNames("code")
	c.SetParamValues(id)

	h = NewURLHandler(suite.echo, suite.repository, &utils.Config{BaseURL: ""})
	return
}

func (suite *URLHandlerTestSuite) prepareShortenTest(url string, rec *httptest.ResponseRecorder) (c echo.Context, h *URLHandler) {
	req := httptest.NewRequest(http.MethodPost, "/shorten", buildPayload(url))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	c = suite.echo.NewContext(req, rec)
	h = NewURLHandler(suite.echo, suite.repository, &utils.Config{BaseURL: ""})
	return
}

func buildPayload(url string) *strings.Reader {
	if url == "" {
		// build invalid payload
		return strings.NewReader("invalid")
	}
	return strings.NewReader(fmt.Sprintf("{%q:%q}", "url", url))
}
