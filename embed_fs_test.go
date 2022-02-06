//go:build test

package main

import (
	"github.com/schaermu/hpfr-shortener/internal/handlers"
	"net/http"
	"testing/fstest"
)

func getFileSystem() http.FileSystem {
	fsys := fstest.MapFS{
		"index.html": {},
	}

	var fs = http.FS(fsys)
	handlers.StaticFS = fs

	return handlers.StaticFS
}
