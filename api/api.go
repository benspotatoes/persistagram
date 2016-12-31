package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/stacktic/dropbox"
	"github.com/tanookiben/go-instagram/instagram"
	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"
)

type Router struct {
	Config    *config
	Dropbox   *dropbox.Dropbox
	Instagram *instagram.Client
}

type errorResponse struct {
	Code   int
	Status string
}

var (
	ErrConfigNotReady = errors.New("Config has not been initialized")
)

func NewRouter() {
	var conf config
	conf.DropboxClientID = os.Getenv("DB_CLIENT_ID")
	conf.DropboxClientSecret = os.Getenv("DB_CLIENT_SECRET")
	conf.DropboxAccessToken = os.Getenv("DB_ACCESS_TOKEN")
	conf.LikedTxtPath = os.Getenv("DB_LIKED_TXT_PATH")

	if !(&conf).Ready() {
		log.Fatal(ErrConfigNotReady)
	}

	db := dropbox.NewDropbox()
	db.SetAppInfo(conf.DropboxClientID, conf.DropboxClientSecret)
	db.SetAccessToken(conf.DropboxAccessToken)

	router := Router{
		Config:  &conf,
		Dropbox: db,
	}

	goji.Get("/running", router.healthCheck)
	goji.Get("/save_liked", router.saveLiked)

	return
}

func (T *Router) healthCheck(c web.C, w http.ResponseWriter, r *http.Request) {
	err := json.NewEncoder(w).Encode(map[string]string{"status": "Running", "code": "200"})
	if err != nil {
		T.serveError(w, r, err)
	}
}

func (T *Router) serveError(w http.ResponseWriter, r *http.Request, e error) {
	log.Printf("returning error: %q", e)
	response := errorResponse{
		Code:   500,
		Status: e.Error(),
	}
	w.WriteHeader(response.Code)
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Fatal(err)
	}
}
