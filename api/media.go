package api

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/benspotatoes/persistagram/backend"
	"github.com/zenazn/goji/web"
)

var (
	filenameRgx = regexp.MustCompile(`\w+\.(mp4|jpg|png)`)

	client = &http.Client{Timeout: 1 * time.Minute}

	mp4Re = regexp.MustCompile(`og:video:secure_url" content="(https:\/\/.*\.mp4)" `)
	jpgRe = regexp.MustCompile(`og:image" content="(https:\/\/.*\.jpg)`)
	ctnRe = regexp.MustCompile(`og:description" content=".* \(@(.*)\) on Instagram`)
	mdaRe = regexp.MustCompile(`meta name="medium" content="(.*)"`)

	sharedDataRe = regexp.MustCompile(`sharedData = (\{.*\})`)
)

const (
	retryCount    = 5
	retryDuration = 30 * time.Second
)

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

func (rt *Router) saveLiked(c web.C, w http.ResponseWriter, r *http.Request) {
	err := rt.Dropbox.DownloadToFile(rt.Config.LikedTxtPath, "/tmp/liked.txt.tmp", "")
	if err != nil {
		rt.serveError(w, r, err)
		return
	}

	raw, err := ioutil.ReadFile("/tmp/liked.txt.tmp")
	if err != nil {
		rt.serveError(w, r, err)
		return
	}
	defer os.Remove("/tmp/liked.txt.tmp")

	liked := strings.Split(string(raw), "\n")
	for _, ig := range liked {
		if ig == "" {
			continue
		}

		log.Println(ig)
		resp, err := client.Get(ig)
		if err != nil {
			log.Printf("error getting url %s: %q\n", ig, err)
			continue
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("error reading body for url %s: %q\n", ig, err)
			continue
		}
		sbody := string(body)

		meta := ctnRe.FindStringSubmatch(sbody)
		if len(meta) != 2 {
			log.Println("unable to get username")
			continue
		}
		username := meta[1]

		src := make(map[string]bool)

		sharedData := sharedDataRe.FindStringSubmatch(sbody)
		if len(sharedData) == 2 {
			blob := sharedData[1]
			data := SharedData{}
			err := json.Unmarshal([]byte(blob), &data)
			if err != nil {
				log.Printf("unable to unmarshal shared data: %q", err)
				continue
			}
			entryData := data.EntryData
			if entryData == nil {
				log.Println("empty entry data")
				continue
			}
			postPage := entryData.PostPage
			if len(postPage) == 0 {
				log.Println("empty post page")
				continue
			}
			for _, post := range data.EntryData.PostPage {
				graph := post.Graphql
				if graph == nil {
					log.Println("empty graphql")
					continue
				}
				media := graph.ShortcodeMedia
				if media == nil {
					log.Println("empty shortcode media")
					continue
				}
				sidecar := media.EdgeSidecarToChildren
				if sidecar == nil {
					log.Println("empty edge sidecar to children")
					continue
				}
				edges := sidecar.Edges
				if len(edges) == 0 {
					log.Println("empty sidecar edges")
					continue
				}
				for _, edge := range edges {
					node := edge.Node
					if node == nil {
						log.Println("empty edge node")
						continue
					}
					displayURL := node.DisplayUrl
					if displayURL == "" {
						log.Println("empty display url")
						continue
					}
					src[displayURL] = true
				}
			}
		}

		metadata := mdaRe.FindStringSubmatch(sbody)
		if len(metadata) != 2 {
			log.Println("unable to get medium type")
			continue
		}
		mediaType := metadata[1]

		switch mediaType {
		case "image":
			img := jpgRe.FindStringSubmatch(sbody)
			if len(img) != 2 {
				log.Println("unable to get image source")
				continue
			}
			src[img[1]] = true
		case "video":
			vid := mp4Re.FindStringSubmatch(sbody)
			if len(vid) != 2 {
				log.Println("unable to get video source")
				continue
			}
			src[vid[1]] = true
		}

		for s, _ := range src {
			data := backend.InstagramMetadata{
				Author:   username,
				Source:   s,
				Filename: getMediaFilename(s),
			}
			log.Println(data)

			go func() {
				retry := 0
				for {
					err := backend.SaveMedia(data, rt.Dropbox)
					if err != nil {
						if retry >= retryCount {
							log.Printf("error saving media for url %s: %q\n", data.Source, err)
							return
						}
						time.Sleep(time.Duration(retry) * retryDuration)
					}
					return
				}
			}()
		}
	}
	_, err = rt.Dropbox.Delete(rt.Config.LikedTxtPath)
	if err != nil {
		rt.serveError(w, r, err)
	}
}

func getMediaFilename(source string) string {
	return filenameRgx.FindString(source)
}
