package instagram

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
)

type Metadata struct {
	Author string
	Links  map[string]string
}

type SharedData struct {
	EntryData *EntryData `json:"entry_data"`
}
type EntryData struct {
	PostPage []*PostPage `json:"PostPage"`
}
type PostPage struct {
	// Media *Media `json:"media"`
	Graphql *Graphql `json:"graphql"`
}
type Graphql struct {
	ShortcodeMedia *ShortcodeMedia `json:"shortcode_media"`
}
type ShortcodeMedia struct {
	EdgeSidecarToChildren *EdgeSidecarToChildren `json:"edge_sidecar_to_children"`
}

// type Media struct {
// 	EdgeSidecarToChildren *EdgeSidecarToChildren `json:"edge_sidecar_to_children"`
// }

type EdgeSidecarToChildren struct {
	Edges []*Edge `json:"edges"`
}
type Edge struct {
	Node *Node `json:"node"`
}
type Node struct {
	DisplayUrl string `json:"display_url"`
	VideoURL   string `json:"video_url"`
}

var (
	filenameRe = regexp.MustCompile(`\w+\.(mp4|jpg|png)`)

	videoRe = regexp.MustCompile(`og:video:secure_url" content="(https:\/\/.*\.mp4)" `)
	imageRe = regexp.MustCompile(`og:image" content="(https:\/\/.*\.jpg)`)

	contentRe  = regexp.MustCompile(`og:description" content=".* \(@(.*)\) on Instagram`)
	content2Re = regexp.MustCompile(`og:description" content=".* @(.*) on Instagram`)

	mediumRe = regexp.MustCompile(`meta name="medium" content="(.*)"`)

	sharedDataRe = regexp.MustCompile(`sharedData = (\{.*\})`)
)

const (
	imageType = "image"
	videoType = "video"
)

func Parse(source string, debug bool) (*Metadata, error) {
	content := contentRe.FindStringSubmatch(source)
	if len(content) != 2 {
		content = content2Re.FindStringSubmatch(source)
		if len(content) != 2 {
			return nil, fmt.Errorf("unable to get username for source %s", source)
		}
	}
	username := content[1]

	links := map[string]string{}

	sharedData := sharedDataRe.FindStringSubmatch(source)
	if len(sharedData) == 2 {
		data := SharedData{}
		err := json.Unmarshal([]byte(sharedData[1]), &data)
		if err != nil {
			if debug {
				log.Printf("unable to get shared data for source %s", source)
			}
		}

		if entryData := data.EntryData; entryData == nil {
			if debug {
				log.Printf("unable to get entry data for source %s", source)
			}
		}

		if postPage := data.EntryData.PostPage; postPage == nil {
			if debug {
				log.Printf("unable to get post page for source %s", source)
			}
		}

		for _, post := range data.EntryData.PostPage {
			if graph := post.Graphql; graph == nil {
				if debug {
					log.Printf("unable to get graphql for source %s", source)
				}
				continue
			}

			if media := post.Graphql.ShortcodeMedia; media == nil {
				if debug {
					log.Printf("unable to get shortcode media for source %s", source)
				}
				continue
			}

			if sidecar := post.Graphql.ShortcodeMedia.EdgeSidecarToChildren; sidecar == nil {
				if debug {
					log.Printf("unable to get edge sidecar to children for source %s", source)
				}
				continue
			}

			if edges := post.Graphql.ShortcodeMedia.EdgeSidecarToChildren.Edges; len(edges) == 0 {
				if debug {
					log.Printf("unable to get edges for source %s", source)
				}
				continue
			}

			for _, edge := range post.Graphql.ShortcodeMedia.EdgeSidecarToChildren.Edges {
				if node := edge.Node; node == nil {
					if debug {
						log.Printf("unable to get node for source %s", source)
					}
					continue
				}

				displayURL := edge.Node.DisplayUrl
				if displayURL != "" {
					links[displayURL] = filename(displayURL)
				}

				videoURL := edge.Node.VideoURL
				if videoURL != "" {
					links[videoURL] = filename(videoURL)
				}
			}
		}
	}

	metadata := mediumRe.FindStringSubmatch(source)
	if len(metadata) != 2 {
		return nil, fmt.Errorf("unable to get medium for source %s", source)
	}
	mediaType := metadata[1]

	switch mediaType {
	case imageType:
		img := imageRe.FindStringSubmatch(source)
		if len(img) != 2 {
			return nil, fmt.Errorf("unable to get image for source %s", source)
		}
		links[img[1]] = filename(img[1])
	case videoType:
		vid := videoRe.FindStringSubmatch(source)
		if len(vid) != 2 {
			return nil, fmt.Errorf("unable to get video for source %s", source)
		}
		links[vid[1]] = filename(vid[1])
	}
	return &Metadata{Author: username, Links: links}, nil
}

func filename(file string) string {
	return filenameRe.FindString(file)
}
