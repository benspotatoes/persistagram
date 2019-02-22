package backend

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"
)

const (
	randSleep = 5 // time.Second
)

func (b *backendImpl) Save(link string) {
	parsed := b.parse([]string{link})
	for _, data := range parsed {
		log.Printf("Saving link %s (%s: %s)", data.path, data.author, data.filename)
		if err := b.save(data); err != nil {
			log.Printf("Unable to save link %s: %s", data.path, err)
		}
	}
}

func (b *backendImpl) save(data *metadata) error {
	// Force re-save
	// if b.exists(data.remoteFilename()) {
	// 	// Skip files that have already been saved
	// 	return nil
	// }

	// Fetch data (retrying up to three times)
	var res *http.Response
	var err error
	res, err = client.Get(data.path)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("error reading response body: %s", err)
		}
		return fmt.Errorf("error fetching data: received code %d with response body '%s'", res.StatusCode, string(body))
	}

	// Upload data (retrying up to three times)
	for range attempts {
		err = b.db.Upload(data.remoteFilename(), res.Body)
		if err == nil {
			break
		}
		time.Sleep(time.Duration(rand.Intn(randSleep)) * time.Second)
	}

	if err != nil {
		return err
	}
	return nil
}

func (b *backendImpl) exists(file string) bool {
	_, err := b.db.Stat(file)
	return err == nil
}
