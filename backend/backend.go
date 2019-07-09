package backend

import (
	"context"
	"log"
	"os"

	"cloud.google.com/go/storage"
	dropbox "github.com/tj/go-dropbox"
	dropy "github.com/tj/go-dropy"
)

type Backend interface {
	Poll()
	Save(link string)
}

type backendImpl struct {
	db        *dropy.Client
	bucket    *storage.BucketHandle
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
	backend := &backendImpl{db, nil, likedFile, saveDir}
	if name := os.Getenv("GCS_BUCKET"); name != "" {
		gcs, err := storage.NewClient(context.Background())
		if err != nil {
			log.Fatalf("Unable to initialize Storage client %s", err)
		}
		bucket := gcs.Bucket(name)
		backend.bucket = bucket
	}
	return backend
}
