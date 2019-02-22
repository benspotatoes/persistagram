package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/benspotatoes/persistagram/backend"
)

func main() {
	backend := backend.NewBackend()
	backend.Poll()

	kill := make(chan os.Signal, 1)
	signal.Notify(kill, syscall.SIGINT, syscall.SIGTERM)

	cron := time.NewTicker(4 * time.Hour).C

	for {
		select {
		case <-cron:
			backend.Poll()
		case <-kill:
			return
		}
	}
}
