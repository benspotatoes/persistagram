package backend

import (
	"fmt"

	dropy "github.com/tj/go-dropy"
)

type InstagramMetadata struct {
	ID       string `json:"id"`
	Author   string `json:"author"`
	Source   string `json:"source"`
	Filename string `json:"filename"`
}

type Backend interface {
	Exists(data InstagramMetadata) bool
	Save(data InstagramMetadata) error
}

type backendImpl struct {
	db *dropy.Client
}

func NewBackend(db *dropy.Client) Backend {
	return &backendImpl{db}
}

func (data *InstagramMetadata) localFilename() string {
	return fmt.Sprintf("/tmp/%s", data.Filename)
}

func (data *InstagramMetadata) remoteFilename() string {
	return fmt.Sprintf("/%s/%s", data.Author, data.Filename)
}
