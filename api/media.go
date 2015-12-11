package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"

	"github.com/benspotatoes/fulsome-corgi/backend"
	"github.com/tanookiben/go-instagram/instagram"
	"github.com/zenazn/goji/web"
)

var (
	ErrMediaMissingResolutions = errors.New("Media missing resolutions")
	ErrUnhandledMediaType      = errors.New("Unhandled media type")

	filenameRgx = regexp.MustCompile("\\w+\\.(mp4|jpg|png)")
)

type recentlyLiked struct {
	RecentlyLiked []backend.InstagramMetadata `json:"recentlyLiked"`
}

type saveMediaResp struct {
	Unsaved []backend.InstagramMetadata `json:"unsavedMedia"`
	Count   int                         `json:"numSaved"`
}

func (T *Router) showRecentlyLiked(c web.C, w http.ResponseWriter, r *http.Request) {
	usersService := T.Instagram.Users

	var count int
	count, err := strconv.Atoi(c.URLParams["count"])
	if err != nil {
		count = 10
	}
	opt := &instagram.Parameters{
		Count: uint64(count),
	}

	likedMedia, _, err := usersService.LikedMedia(opt)
	if err != nil {
		T.serveError(w, r, err)
	}

	var liked []backend.InstagramMetadata
	// liked := make([]backend.InstagramMetadata, count)
	for _, media := range likedMedia {
		mediaSource, err := getMediaSource(media)
		if err != nil {
			log.Println(err)
		}
		metadata := backend.InstagramMetadata{
			ID:       media.ID,
			Author:   media.User.Username,
			Source:   mediaSource,
			Filename: getMediaFilename(mediaSource),
		}
		metadata.Saved = backend.MediaSaved(metadata, T.Dropbox)
		liked = append(liked, metadata)
	}

	w.WriteHeader(200)
	err = json.NewEncoder(w).Encode(recentlyLiked{
		RecentlyLiked: liked,
	})
	if err != nil {
		T.serveError(w, r, err)
	}
}

// TODO - add max allowed number of media to save per run?
func (T *Router) saveInstagramMedia(c web.C, w http.ResponseWriter, r *http.Request) {
	var err error

	usersService := T.Instagram.Users
	opt := &instagram.Parameters{
		Count: uint64(T.Config.igPaginationCount),
	}

	var lastSavedUpdated bool
	var breakSave bool
	var countSaved int

	refLastSavedMediaID := T.Config.InstagramLastSavedMediaID
	var unsavedMedia []backend.InstagramMetadata
	for !breakSave {
		likedMedia, pagination, err := usersService.LikedMedia(opt)
		if err != nil {
			T.serveError(w, r, err)
		}

		for _, media := range likedMedia {
			log.Println(fmt.Sprintf("Parsing media ID '%s'", media.ID))
			var saved bool
			// Set the new "InstagramLastSavedMediaID" value: we want the first media
			// ID processed to be come the future reference of where we should stop
			// processing
			if !lastSavedUpdated {
				T.Config.InstagramLastSavedMediaID = media.ID
				lastSavedUpdated = true
			}

			// Check to see if we have reached the end of un-processed liked media
			if media.ID == refLastSavedMediaID {
				breakSave = true
			}

			// Continue processing liked media if we haven't reached the stopping
			// point of the last run
			if !breakSave {
				mediaSource, err := getMediaSource(media)
				if err != nil {
					log.Println(err)
				}
				metadata := backend.InstagramMetadata{
					ID:       media.ID,
					Author:   media.User.Username,
					Source:   mediaSource,
					Filename: getMediaFilename(mediaSource),
				}

				// Check to make sure we haven't already downloaded this item
				if !backend.MediaSaved(metadata, T.Dropbox) {
					err = backend.SaveMedia(metadata, T.Dropbox)
					if err != nil {
						log.Println(err)
					} else {
						saved = true
						countSaved += 1
					}
				}

				// Update pagination struct to return the next batch of unprocessed
				// liked media
				opt.MaxID = pagination.NextMaxLikeID

				// Keep track of unsaved media for response
				if !saved {
					metadata.Saved = false
					unsavedMedia = append(unsavedMedia, metadata)
				}
			}
		}
	}

	err = ioutil.WriteFile(T.Config.InstagramLastSavedPath, []byte(T.Config.InstagramLastSavedMediaID), 0644)
	if err != nil {
		T.serveError(w, r, err)
	}

	w.WriteHeader(200)
	err = json.NewEncoder(w).Encode(saveMediaResp{
		Unsaved: unsavedMedia,
		Count:   countSaved,
	})
	if err != nil {
		T.serveError(w, r, err)
	}
}

func getMediaSource(m instagram.Media) (string, error) {
	var source string
	var err error
	switch m.Type {
	case "image":
		source, err = getImageSource(m.Images)
		if err != nil {
			return source, fmt.Errorf("%s: %s", err.Error(), m.ID)
		}
	case "video":
		source, err = getVideoSource(m.Videos)
		if err != nil {
			return source, fmt.Errorf("%s: %s", err.Error(), m.ID)
		}
	default:
		log.Println(fmt.Errorf("%s: %s::%s", ErrUnhandledMediaType, m.ID, m.Type))
	}
	return source, nil
}

func getMediaFilename(source string) string {
	return filenameRgx.FindString(source)
}

func getVideoSource(mv *instagram.MediaVideos) (string, error) {
	var source string
	if low := mv.LowResolution; low != nil {
		source = low.URL
	}
	if standard := mv.StandardResolution; standard != nil {
		source = standard.URL
	}
	if source == "" {
		return source, ErrMediaMissingResolutions
	}

	return source, nil
}

func getImageSource(mi *instagram.MediaImages) (string, error) {
	var source string
	if low := mi.LowResolution; low != nil {
		source = low.URL
	}
	if standard := mi.StandardResolution; standard != nil {
		source = standard.URL
	}
	if source == "" {
		return source, ErrMediaMissingResolutions
	}

	return source, nil
}
