package api

import (
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/benspotatoes/persistagram/backend"
	"github.com/benspotatoes/persistagram/internal/instagram"
)

var (
	client = &http.Client{Timeout: 1 * time.Minute}
)

const (
	likedTxt = "/liked.txt"

	retryCount    = 5
	retryDuration = 30 * time.Second
	randSleep     = 5000
)

func (t *Router) save(w http.ResponseWriter, r *http.Request) {
	liked, err := t.liked()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("unable to get liked: %q", err)
		return
	}

	t.process(liked)

	if err := t.clean(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("unable to clean: %q", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (t *Router) liked() ([]string, error) {
	data, err := t.Dropbox.Download(likedTxt)
	if err != nil {
		return []string{}, err
	}
	raw, err := ioutil.ReadAll(data)
	if err != nil {
		return []string{}, err
	}
	return strings.Split(string(raw), "\n"), nil
}

func (t *Router) process(liked []string) {
	for _, link := range liked {
		if link == "" {
			continue
		}

		res, err := client.Get(link)
		if err != nil {
			log.Printf("unable to get link %s: %q", link, err)
			continue
		}
		defer res.Body.Close()

		raw, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Printf("unable to read body for link %s: %q", link, err)
			continue
		}

		info, err := instagram.Parse(string(raw))
		if err != nil {
			log.Printf("unable to parse body for link %s: %q", link, err)
			continue
		}

		for path, filename := range info.Links {
			data := backend.InstagramMetadata{
				Author:   info.Author,
				Source:   path,
				Filename: filename,
			}
			if t.Backend.Exists(data) {
				log.Printf("file exists for link %s", link)
				continue
			}

			go func() {
				time.Sleep(time.Duration(rand.Intn(randSleep)) * time.Millisecond)
				try := 0
				for {
					if err := t.Backend.Save(data); err != nil {
						try = try + 1
						if try >= retryCount {
							log.Printf("error saving source for link %s: %q", link, err)
							return
						}
						time.Sleep(time.Duration(try) * retryDuration)
						continue
					}
					return
				}
			}()
		}
	}
	return
}

func (t *Router) clean() error {
	return t.Dropbox.Delete(likedTxt)
}
