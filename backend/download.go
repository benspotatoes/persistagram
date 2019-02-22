package backend

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func (b *backendImpl) download(data *metadata) error {
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

	author := data.safeAuthor()
	saveDirectory := fmt.Sprintf("%s/%s", b.saveDir, author)
	if err := os.MkdirAll(saveDirectory, os.ModePerm); err != nil {
		return fmt.Errorf("error creating directory (%s): %s", author, err)
	}

	savePath := fmt.Sprintf("%s/%s/%s", b.saveDir, data.safeAuthor(), data.filename)
	if err := ioutil.WriteFile(savePath, body, os.ModePerm); err != nil {
		return fmt.Errorf("error writing file: %s", err)
	}
	return nil
}
