package backend

import (
	"strings"
)

type metadata struct {
	author   string
	filename string
	path     string
}

func (data *metadata) safeAuthor() string {
	author := data.author
	author = strings.Replace(author, ".", "_", -1)
	author = strings.Replace(author, "-", "_", -1)
	return author
}
