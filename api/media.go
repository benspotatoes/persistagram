package api

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/benspotatoes/fulsome-corgi/backend"
	"github.com/zenazn/goji/web"
)

var (
	filenameRgx = regexp.MustCompile(`\w+\.(mp4|jpg|png)`)

	client = &http.Client{Timeout: 1 * time.Minute}

	mp4Re = regexp.MustCompile(`og:video:secure_url" content="(https:\/\/.*\.mp4)" `)
	jpgRe = regexp.MustCompile(`og:image" content="(https:\/\/.*\.jpg)`)
	ctnRe = regexp.MustCompile(`og:description" content="See this Instagram (\w+) by @(\w+)`)
)

func (rt *Router) saveLiked(c web.C, w http.ResponseWriter, r *http.Request) {
	err := rt.Dropbox.DownloadToFile(rt.Config.LikedTxtPath, "liked.txt.tmp", "")
	if err != nil {
		rt.serveError(w, r, err)
		return
	}

	raw, err := ioutil.ReadFile("liked.txt.tmp")
	if err != nil {
		rt.serveError(w, r, err)
		return
	}
	defer os.Remove("liked.txt.tmp")

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
		if len(meta) != 3 {
			log.Println("unable to get metadata")
			continue
		}

		username := meta[2]
		mediaType := meta[1]

		var src string
		switch mediaType {
		case "photo":
			img := jpgRe.FindStringSubmatch(sbody)
			if len(img) != 2 {
				log.Println("unable to get photo source")
				continue
			}
			src = img[1]
		case "video":
			vid := mp4Re.FindStringSubmatch(sbody)
			if len(vid) != 2 {
				log.Println("unable to get video source")
				continue
			}
			src = vid[1]
		default:
			continue
		}

		data := backend.InstagramMetadata{
			Author:   username,
			Source:   src,
			Filename: getMediaFilename(src),
		}
		log.Println(data)

		go func() {
			err := backend.SaveMedia(data, rt.Dropbox)
			if err != nil {
				log.Printf("error saving media for url %s: %q\n", data.Source, err)
			}
		}()
	}
	_, err = rt.Dropbox.Delete(rt.Config.LikedTxtPath)
	if err != nil {
		rt.serveError(w, r, err)
	}
}

func getMediaFilename(source string) string {
	return filenameRgx.FindString(source)
}
