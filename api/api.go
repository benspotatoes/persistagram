package api

import (
	"net/http"

	"goji.io"
	"goji.io/pat"

	"github.com/benspotatoes/persistagram/backend"
	dropy "github.com/tj/go-dropy"
)

type Router struct {
	*goji.Mux
	Dropbox *dropy.Client
	Backend backend.Backend
}

type errorResponse struct {
	Code   int
	Status string
}

func NewRouter(b backend.Backend, db *dropy.Client) *Router {
	mux := goji.NewMux()

	router := &Router{
		Mux:     mux,
		Dropbox: db,
		Backend: b,
	}

	router.HandleFunc(pat.Get("/running"), router.healthCheck)
	router.HandleFunc(pat.Get("/save"), router.save)

	return router
}

func (t *Router) healthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
}
