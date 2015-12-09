package api

import (
	"encoding/json"
	"errors"
	"io/ioutil"
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

	conf.InstagramClientID = os.Getenv("IG_CLIENT_ID")
	conf.InstagramClientSecret = os.Getenv("IG_CLIENT_SECRET")
	conf.InstagramRedirectURL = os.Getenv("IG_REDIRECT_URL")
	conf.InstagramAccessTokenPath = os.Getenv("IG_ACCESS_TOKEN_PATH")
	conf.InstagramLastSavedPath = os.Getenv("IG_LAST_SAVED_PATH")
	conf.InstagramLastSavedMediaID = os.Getenv("IG_LAST_SAVED_MEDIA_ID")
	conf.DropboxClientID = os.Getenv("DB_CLIENT_ID")
	conf.DropboxClientSecret = os.Getenv("DB_CLIENT_SECRET")
	conf.DropboxAccessToken = os.Getenv("DB_ACCESS_TOKEN")

	if !(&conf).Ready() {
		log.Fatal(ErrConfigNotReady)
	}

	instagram := instagram.NewClient(nil)
	instagram.ClientID = conf.InstagramClientID
	instagram.ClientSecret = conf.InstagramClientSecret
	tkn, err := ioutil.ReadFile(conf.InstagramAccessTokenPath)
	if err != nil {
		log.Println(err)
	}
	instagramToken := string(tkn)
	if instagramToken != "" {
		conf.InstagramAccessToken = instagramToken
		instagram.AccessToken = instagramToken
	}

	saved, err := ioutil.ReadFile(conf.InstagramLastSavedPath)
	if err != nil {
		log.Println(err)
	}
	lastSavedMediaID := string(saved)
	if lastSavedMediaID != "" {
		conf.InstagramLastSavedMediaID = lastSavedMediaID
	}

	db := dropbox.NewDropbox()
	db.SetAppInfo(conf.DropboxClientID, conf.DropboxClientSecret)
	db.SetAccessToken(conf.DropboxAccessToken)

	router := Router{
		Config:    &conf,
		Instagram: instagram,
		Dropbox:   db,
	}

	goji.Get("/running", router.healthCheck)
	goji.Get("/fetch_instagram_token", router.fetchInstagramToken)
	goji.Get("/save_instagram_media", router.saveInstagramMedia)

	return
}

func (T *Router) healthCheck(c web.C, w http.ResponseWriter, r *http.Request) {
	err := json.NewEncoder(w).Encode(map[string]string{"status": "Running", "code": "200"})
	if err != nil {
		T.serveError(w, r, err)
	}
}

func (T *Router) serveError(w http.ResponseWriter, r *http.Request, e error) {
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
