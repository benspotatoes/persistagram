package backend

type InstagramMetadata struct {
	ID       string `json:"id"`
	Author   string `json:"author"`
	Source   string `json:"source"`
	Filename string `json:"filename"`
	Saved    bool   `json:"saved"`
}
