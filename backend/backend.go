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
	db        *dropy.Client
	likedFile string
	saveDir   string
}

var (
	attempts = []int{1, 2, 3}
)

func NewBackend() Backend {
	db := dropy.New(dropbox.New(dropbox.NewConfig(os.Getenv("DB_ACCESS_TOKEN"))))
	likedFile := os.Getenv("LIKED_FILE")
	if likedFile == "" {
		likedFile = "/liked.txt"
	}
	saveDir := os.Getenv("SAVE_DIRECTORY")
	if saveDir == "" {
		saveDir = "/opt/persistagram/data"
	}
	return &backendImpl{db, likedFile, saveDir}
}
