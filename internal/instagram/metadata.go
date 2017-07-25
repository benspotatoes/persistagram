package instagram

import (
	"encoding/json"
	"errors"
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

var (
	errUsername              = errors.New("unable to get username")
	errSharedData            = errors.New("unable to get shared data")
	errEntryData             = errors.New("unable to get entry data")
	errPostPage              = errors.New("unable to get post page")
	errGraphql               = errors.New("unable to get graphql")
	errShortcodeMedia        = errors.New("unable to get shortcode media")
	errEdgeSidecarToChildren = errors.New("unable to get edge sidecar to children")
	errEdges                 = errors.New("unable to get edges")
	errNode                  = errors.New("unable to get node")
	errDisplayURL            = errors.New("unable to get display url")
	errMedium                = errors.New("unable to get medium")

	errImage = errors.New("unable to get image source")
	errVideo = errors.New("unable to get video source")
)

const (
	imageType = "image"
	videoType = "video"
)

func Parse(source string) (*Metadata, error) {
	content := contentRe.FindStringSubmatch(source)
	if len(content) != 2 {
		content = content2Re.FindStringSubmatch(source)
		if len(content) != 2 {
			return nil, errUsername
		}
	}
	username := content[1]

	links := map[string]string{}

	sharedData := sharedDataRe.FindStringSubmatch(source)
	if len(sharedData) == 2 {
		data := SharedData{}
		err := json.Unmarshal([]byte(sharedData[1]), &data)
		if err != nil {
			log.Printf(errSharedData.Error())
		}

		if entryData := data.EntryData; entryData == nil {
			log.Printf(errEntryData.Error())
		}

		if postPage := data.EntryData.PostPage; postPage == nil {
			log.Printf(errPostPage.Error())
		}

		for _, post := range data.EntryData.PostPage {
			if graph := post.Graphql; graph == nil {
				log.Printf(errGraphql.Error())
				continue
			}

			if media := post.Graphql.ShortcodeMedia; media == nil {
				log.Printf(errShortcodeMedia.Error())
				continue
			}

			if sidecar := post.Graphql.ShortcodeMedia.EdgeSidecarToChildren; sidecar == nil {
				log.Printf(errEdgeSidecarToChildren.Error())
				continue
			}

			if edges := post.Graphql.ShortcodeMedia.EdgeSidecarToChildren.Edges; len(edges) == 0 {
				log.Printf(errEdges.Error())
				continue
			}

			for _, edge := range post.Graphql.ShortcodeMedia.EdgeSidecarToChildren.Edges {
				if node := edge.Node; node == nil {
					log.Printf(errNode.Error())
					continue
				}

				displayURL := edge.Node.DisplayUrl
				if displayURL == "" {
					log.Printf(errDisplayURL.Error())
					continue
				}

				links[displayURL] = filename(displayURL)
			}
		}
	}

	metadata := mediumRe.FindStringSubmatch(source)
	if len(metadata) != 2 {
		return nil, errMedium
	}
	mediaType := metadata[1]

	switch mediaType {
	case imageType:
		img := imageRe.FindStringSubmatch(source)
		if len(img) != 2 {
			return nil, errImage
		}
		links[img[1]] = filename(img[1])
	case videoType:
		vid := videoRe.FindStringSubmatch(source)
		if len(vid) != 2 {
			return nil, errVideo
		}
		links[vid[1]] = filename(vid[1])
	}
	return &Metadata{Author: username, Links: links}, nil
}

func filename(file string) string {
	return filenameRe.FindString(file)
}
