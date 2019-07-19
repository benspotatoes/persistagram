package backend

import (
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/benspotatoes/persistagram/clients/instagram"
)

var (
	client = &http.Client{Timeout: 30 * time.Second}
)

// parse takes a list of Instagram links and returns two lists: one list of
// parsed Instagram metadata and one list of links that could not be parsed
// successfully
func (b *backendImpl) parse(liked []string) ([]*metadata, []string) {
	var parsed []*metadata
	var failed []string
	for _, link := range liked {
		// Skip empty links
		if link == "" {
			continue
		}

		page, err := readLink(link)
		if err != nil {
			log.Printf("Unable to read link %s: %s\n", link, err)
			failed = append(failed, link)
			continue
		}

		data, err := instagram.Parse(page, false)
		if err != nil {
			log.Printf("Unable to parse link %s: %s\n", link, err)
			failed = append(failed, link)
			continue
		}

		for path, filename := range data.Links {
			parsed = append(parsed, &metadata{
				author:   data.Author,
				filename: filename,
				path:     path,
			})
		}
	}
	return parsed, failed
}

// readLink fetches the data for the specified link (retrying up to three
// times) and returns the raw data as a string
func readLink(link string) (string, error) {
	var res *http.Response
	var err error
	for range attempts {
		res, err = client.Get(link)
		if err == nil {
			break
		}
	}
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	raw, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	return string(raw), nil
}
