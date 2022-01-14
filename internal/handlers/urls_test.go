package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

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
	mongoServer *memongo.Server
	repository  *repositories.URLRepository
}

func (suite *URLHandlerTestSuite) SetupSuite() {
	mongoServer, err := memongo.Start("5.0.5")
	if err != nil {
		suite.T().Fatal(err)
	}
	store := data.NewDatastore(mongoServer.URI(), memongo.RandomDatabase(), logger)
	suite.mongoServer = mongoServer
	suite.repository = repositories.NewURLRepository(store)
}

func (suite *URLHandlerTestSuite) TearDownSuite() {
	suite.mongoServer.Stop()
}

func TestURLHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(URLHandlerTestSuite))
}

func (suite *URLHandlerTestSuite) TestShorten() {
	for scenario, fn := range map[string]func(){
		"invalid url": suite.testInvalidURL,
	} {
		suite.Run(scenario, fn)
	}
}

func (suite *URLHandlerTestSuite) testInvalidURL() {
	// arrange
	var url = "bogus_url_input"
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/shorten", buildPayload(url))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)
	h := NewURLHandler(e, suite.repository, &utils.Config{BaseURL: ""})

	// act
	if assert.NoError(suite.T(), h.Shorten(c)) {
	}
}

func buildPayload(url string) *strings.Reader {
	return strings.NewReader((fmt.Sprintf("{%q:%q}", "url", url)))
}
