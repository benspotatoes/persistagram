package backend

import (
	"fmt"
	"strings"
)

type metadata struct {
	author   string
	filename string
	path     string
}

func (data *metadata) localFilename() string {
	return fmt.Sprintf("/tmp/%s", data.filename)
}

func (data *metadata) remoteFilename() string {
	return fmt.Sprintf("/%s/%s", data.safeAuthor(), data.filename)
}

func (data *metadata) safeAuthor() string {
	author := data.author
	author = strings.Replace(author, ".", "_", -1)
	author = strings.Replace(author, "-", "_", -1)
	return author
}
