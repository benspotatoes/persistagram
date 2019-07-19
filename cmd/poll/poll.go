package main

import (
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/benspotatoes/persistagram/backend"
)

func main() {
	backend := backend.NewBackend()
	backend.Poll()

	kill := make(chan os.Signal, 1)
	signal.Notify(kill, syscall.SIGINT, syscall.SIGTERM)

	var interval int64
	var err error
	interval, err = strconv.ParseInt(os.Getenv("INTERVAL"), 10, 64)
	if err != nil {
		interval = 4
	}

	cron := time.NewTicker(time.Duration(interval) * time.Hour).C

	for {
		select {
		case <-cron:
			backend.Poll()
		case <-kill:
			return
		}
	}
}
