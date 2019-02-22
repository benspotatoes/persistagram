package backend

import (
	"os"

	dropbox "github.com/tj/go-dropbox"
	dropy "github.com/tj/go-dropy"
)

type Backend interface {
	Process()
	Save(link string)
}

type backendImpl struct {
	db *dropy.Client
}

func NewBackend() Backend {
	db := dropy.New(dropbox.New(dropbox.NewConfig(os.Getenv("DB_ACCESS_TOKEN"))))
	return &backendImpl{db}
}
