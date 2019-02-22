package backend

import (
	"io/ioutil"
	"log"
	"strings"
)

var (
	attempts = []int{1, 2, 3}
)

const (
	likedTxt   = "/liked.txt"
	likedDelim = "\n"
)

func (b *backendImpl) Process() {
	if !b.exists(likedTxt) {
		return
	}

	liked, err := b.get()
	if err != nil {
		log.Printf("Unable to get liked file: %s", err)
		return
	}

	parsed := b.parse(liked)
	for _, data := range parsed {
		go func(data *metadata) {
			if err := b.save(data); err != nil {
				log.Printf("Unable to save link %s: %s", data.path, err)
			}
		}(data)
	}

	if err := b.clean(); err != nil {
		log.Printf("Unable to clean: %s", err)
	}
}

func (b *backendImpl) get() ([]string, error) {
	data, err := b.db.Download(likedTxt)
	if err != nil {
		return []string{}, err
	}
	raw, err := ioutil.ReadAll(data)
	if err != nil {
		return []string{}, err
	}
	return strings.Split(string(raw), likedDelim), nil
}

func (b *backendImpl) clean() error {
	return b.db.Delete(likedTxt)
}
