package main

import (
	"github.com/benspotatoes/fulsome-corgi/api"

	"github.com/zenazn/goji"
)

func main() {
	api.NewRouter()
	goji.Serve()
}
