package utils

import (
	"net/http"
	"testing/fstest"
)

func GetFakeFileSystem() http.FileSystem {
	fsys := fstest.MapFS{
		"index.html": {},
	}
	return http.FS(fsys)
}
