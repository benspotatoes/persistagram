package backend

import (
	"math/rand"
	"net/http"
	"time"
)

const (
	randSleep = 5 // time.Second
)

func (b *backendImpl) save(data *metadata) error {
	if b.exists(data.remoteFilename()) {
		// Skip files that have already been saved
		return nil
	}

	// Fetch data (retrying up to three times)
	var res *http.Response
	var err error
	for range attempts {
		res, err = client.Get(data.path)
		if err == nil {
			break
		}
		time.Sleep(time.Duration(rand.Intn(randSleep)) * time.Second)
	}
	if err != nil {
		return err
	}

	// Upload data (retrying up to three times)
	for range attempts {
		err = b.db.Upload(data.remoteFilename(), res.Body)
		if err == nil {
			break
		}
		time.Sleep(time.Duration(rand.Intn(randSleep)) * time.Second)
	}

	// Make sure we don't close the body before we're done uploading
	defer res.Body.Close()

	if err != nil {
		return err
	}
	return nil
}

func (b *backendImpl) exists(file string) bool {
	_, err := b.db.Stat(file)
	return err == nil
}
