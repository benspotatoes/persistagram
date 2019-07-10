package backend

import (
	"io/ioutil"
	"log"
	"regexp"
	"strings"
)

var (
	notFoundRe = regexp.MustCompile("not_found")
)

func (b *backendImpl) Poll() {
	liked, err := b.get()
	if err != nil {
		if !notFoundRe.MatchString(err.Error()) {
			log.Printf("Unable to get liked file: %s", err)
		}
		return
	}

	parsed := b.parse(liked)
	for _, data := range parsed {
		if err := b.download(data); err != nil {
			log.Printf("Unable to save link %s: %s", data.path, err)
		}
	}

	if err := b.clean(); err != nil {
		log.Printf("Unable to clean: %s", err)
	}
}

func (b *backendImpl) get() ([]string, error) {
	data, err := b.db.Download(b.likedFile)
	if err != nil {
		return []string{}, err
	}
	raw, err := ioutil.ReadAll(data)
	if err != nil {
		return []string{}, err
	}
	return strings.Split(string(raw), "\n"), nil
}

func (b *backendImpl) clean() error {
	return b.db.Delete(b.likedFile)
}
