package backend

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/stacktic/dropbox"
)

func MediaSaved(data InstagramMetadata, client *dropbox.Dropbox) bool {
	_, err := client.Metadata(fmt.Sprintf("%s/%s", data.Author, data.Filename), false, false, "", "", 1)
	if err != nil {
		return false
	}
	return true
}

func SaveMedia(data InstagramMetadata, client *dropbox.Dropbox) error {
	localFilepath := fmt.Sprintf("/tmp/%s", data.Filename)

	// https://github.com/thbar/golang-playground/blob/master/download-files.go
	output, err := os.Create(localFilepath)
	if err != nil {
		return fmt.Errorf("Unable to create file: %s\n%s", data.Filename, err)
	}
	defer output.Close()
	defer os.Remove(localFilepath)

	response, err := http.Get(data.Source)
	if err != nil {
		return fmt.Errorf("Unable to download file: %s\n%s", data.Filename, err)
	}
	defer response.Body.Close()

	_, err = io.Copy(output, response.Body)
	if err != nil {
		return fmt.Errorf("Unable to download file: %s\n%s", data.Filename, err)
	}

	dropboxFilepath := fmt.Sprintf("%s/%s", data.Author, data.Filename)
	if _, err = client.UploadFile(localFilepath, dropboxFilepath, true, ""); err != nil {
		return fmt.Errorf("Unable to upload file to Dropbox: %s\n%s", dropboxFilepath, err)
	}

	return nil
}
