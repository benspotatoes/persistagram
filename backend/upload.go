package backend

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"
)

func (b *backendImpl) upload(data *metadata) error {
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

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %s", err)
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("error fetching data: received code %d with response body '%s'", res.StatusCode, string(body))
	}

	// Upload data (retrying up to three times)
	for range attempts {
		err = b.db.Upload(data.remoteFilename(), bytes.NewReader(body))
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
