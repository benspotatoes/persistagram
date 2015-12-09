package api

import "strconv"

type config struct {
	// Instagram settings
	InstagramClientID     string
	InstagramClientSecret string
	InstagramRedirectURL  string

	// Optional Instagram settings
	InstagramAccessToken      string
	InstagramAccessTokenPath  string
	InstagramLastSavedPath    string
	InstagramLastSavedMediaID string
	InstagramPaginationCount  string

	// Parsed pagination count
	igPaginationCount int

	// Dropbox settings
	DropboxClientID     string
	DropboxClientSecret string
	DropboxAccessToken  string
}

func (c *config) Ready() bool {
	if c.InstagramClientID == "" {
		return false
	}
	if c.InstagramClientSecret == "" {
		return false
	}
	if c.InstagramRedirectURL == "" {
		return false
	}
	if c.InstagramAccessTokenPath == "" {
		c.InstagramAccessTokenPath = "./ig_access_token"
	}
	if c.InstagramLastSavedPath == "" {
		c.InstagramLastSavedPath = "./ig_last_saved"
	}
	if c.InstagramPaginationCount == "" {
		c.igPaginationCount = 5
	} else {
		val, err := strconv.Atoi(c.InstagramPaginationCount)
		if err != nil {
			val = 5
		}
		c.igPaginationCount = val
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
