//go:build !test

package main

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/schaermu/hpfr-shortener/internal/handlers"
)

//go:embed ui/dist/*
var embeddedFiles embed.FS

func getFileSystem() http.FileSystem {
	fsys, err := fs.Sub(embeddedFiles, "ui/dist")
	if err != nil {
		panic(err)
	}

	var fs = http.FS(fsys)
	handlers.StaticFS = fs

	return fs
}
