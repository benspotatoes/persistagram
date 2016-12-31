package main

import (
	"github.com/benspotatoes/persistagram/api"

	"github.com/zenazn/goji"
)

func main() {
	api.NewRouter()
	goji.Serve()
}
