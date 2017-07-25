package backend

import (
	"fmt"
	"net/http"
	"time"
)

var (
	client = &http.Client{Timeout: 1 * time.Minute}
)

func (b *backendImpl) Save(data InstagramMetadata) error {
	res, err := client.Get(data.Source)
	if err != nil {
		return fmt.Errorf("unable to download file %s: %q", data.Source, err)
	}
	defer res.Body.Close()

	if err := b.db.Upload(data.remoteFilename(), res.Body); err != nil {
		return fmt.Errorf("unable save file: %s: %q", data.Source, err)
	}

	return nil
}
