package backend

import (
	"log"
)

const (
	randSleep = 5 // time.Second
)

func (b *backendImpl) Save(link string) {
	parsed := b.parse([]string{link})
	for _, data := range parsed {
		log.Printf("Saving link %s (%s: %s)", data.path, data.author, data.filename)
		if err := b.upload(data); err != nil {
			log.Printf("Unable to save link %s: %s", data.path, err)
		}
	}
}
