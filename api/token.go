package api

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/zenazn/goji/web"

	"golang.org/x/oauth2"
)

func (T *Router) fetchInstagramToken(c web.C, w http.ResponseWriter, r *http.Request) {
	oauthCfg := &oauth2.Config{
		ClientID:     T.Config.InstagramClientID,
		ClientSecret: T.Config.InstagramClientSecret,
		RedirectURL:  T.Config.InstagramRedirectURL,
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://api.instagram.com/oauth/authorize/",
			TokenURL: "https://api.instagram.com/oauth/access_token",
		},
	}

	url := oauthCfg.AuthCodeURL("state", oauth2.AccessTypeOnline)
	log.Println("Visit URL for auth dialog: ", url)

	var code string
	if _, err := fmt.Scan(&code); err != nil {
		T.serveError(w, r, err)
	}

	tkn, err := oauthCfg.Exchange(oauth2.NoContext, code)
	if err != nil {
		T.serveError(w, r, err)
	}

	instagramToken := tkn.AccessToken
	T.Config.InstagramAccessToken = instagramToken
	T.Instagram.AccessToken = instagramToken
	err = ioutil.WriteFile(T.Config.InstagramAccessTokenPath, []byte(instagramToken), 0644)
	if err != nil {
		T.serveError(w, r, err)
	}
}
