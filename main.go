package main

import (
	"net/http"
	"os"

	"github.com/benspotatoes/persistagram/api"
	"github.com/benspotatoes/persistagram/backend"
	dropbox "github.com/tj/go-dropbox"
	dropy "github.com/tj/go-dropy"
)

func main() {
	db := initDropbox()
	backend := backend.NewBackend(db)
	api := api.NewRouter(backend, db)
	http.ListenAndServe("localhost:8000", api)
}

func initDropbox() *dropy.Client {
	return dropy.New(dropbox.New(dropbox.NewConfig(os.Getenv("DB_ACCESS_TOKEN"))))
}
