package cmd

import (
	"fmt"
	"net/http"
	"os"
	"testing"
	"testing/fstest"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/tryvium-travels/memongo"
)

func getTestFs() http.FileSystem {
	fsys := fstest.MapFS{
		"index.html": {},
	}
	return http.FS(fsys)
}

var log = logrus.New()

func TestStartWithoutConfig(t *testing.T) {
	// act
	_, err := NewServer(getTestFs(), log)

	// assert
	if assert.Error(t, err) {
		assert.EqualError(t, err, fmt.Sprintf("env var %q not set", "MONGO_DSN"))
	}
}

func TestStartWithInvalidConfig(t *testing.T) {
	// arrange
	os.Setenv("MONGO_DSN", "foobar")
	os.Setenv("MONGO_DB", "foobar")

	// act
	_, err := NewServer(getTestFs(), log)

	// assert
	if assert.Error(t, err) {
		assert.EqualError(t, err, "error connecting to database: foobar/foobar")
	}
}

func TestStartWithValidConfig(t *testing.T) {
	// arrange
	mVersion := os.Getenv("MEMONGO_VERSION")
	if mVersion == "" {
		mVersion = "5.0.5"
	}
	mongoServer, err := memongo.Start(mVersion)
	if err != nil {
		t.Fatal(err)
	}
	os.Setenv("MONGO_DSN", mongoServer.URI())
	os.Setenv("MONGO_DB", memongo.RandomDatabase())

	// act
	_, err = NewServer(getTestFs(), log)

	// assert
	assert.NoError(t, err)
}
