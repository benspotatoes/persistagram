package api

type config struct {
	// Dropbox settings
	DropboxClientID     string
	DropboxClientSecret string
	DropboxAccessToken  string
	LikedTxtPath        string
}

func (c *config) Ready() bool {
	if c.LikedTxtPath == "" {
		c.LikedTxtPath = "liked.txt"
	}
	if c.DropboxClientID == "" {
		return false
	}
	if c.DropboxClientSecret == "" {
		return false
	}
	if c.DropboxAccessToken == "" {
		return false
	}
	return true
}
