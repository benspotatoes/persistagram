package backend

import (
	"os"

	dropbox "github.com/tj/go-dropbox"
	dropy "github.com/tj/go-dropy"
)

type Backend interface {
	Poll()
	Save(link string)
}

type backendImpl struct {
	db       *dropy.Client
	saveFile string
}

var (
	attempts = []int{1, 2, 3}
)

func NewBackend() Backend {
	db := dropy.New(dropbox.New(dropbox.NewConfig(os.Getenv("DB_ACCESS_TOKEN"))))
	saveFile := os.Getenv("SAVE_FILE")
	if saveFile == "" {
		saveFile = "/save.txt"
	}
	return &backendImpl{db, saveFile}
}
