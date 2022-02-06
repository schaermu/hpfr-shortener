package utils

import (
	"net/http"
	"testing/fstest"
)

func GetFakeFileSystem() http.FileSystem {
	fsys := fstest.MapFS{
		"index.html": &fstest.MapFile{Data: []byte("<html></html>")},
	}
	return http.FS(fsys)
}
