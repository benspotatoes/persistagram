package main

import (
  "io"
  "os"
  "log"
  "fmt"
  "regexp"
  "strings"
  "net/http"
  "io/ioutil"
  "encoding/json"
  "golang.org/x/oauth2"
  "github.com/stacktic/dropbox"
  "github.com/gedex/go-instagram/instagram"
)

type InstagramMedia struct {
  AuthorName string `json:"author_name"`
  MediaId string `json:"media_id"`
}

func get_ig_access_token() oauth2.Token {
  conf := &oauth2.Config{
    ClientID: os.Getenv("IG_CLIENT_ID"),
    ClientSecret: os.Getenv("IG_CLIENT_SECRET"),
    RedirectURL: os.Getenv("IG_REDIRECT_URL"),
    Endpoint: oauth2.Endpoint{
      AuthURL: "https://api.instagram.com/oauth/authorize/",
      TokenURL: "https://api.instagram.com/oauth/access_token",
    },
  }

  url := conf.AuthCodeURL("state", oauth2.AccessTypeOnline)
  fmt.Printf("Visit the URL for the auth dialog: %v\n", url)

  var code string
  if _, err := fmt.Scan(&code); err != nil {
    log.Fatal(err)
  }
  token, err := conf.Exchange(oauth2.NoContext, code)
  if err != nil {
    log.Fatal(err)
  }

  return *token
}

func get_ig_source(src string, tok string) (author, source string) {
  resp, err := http.Get(fmt.Sprintf("http://api.instagram.com/oembed?url=%s", src))
  if err != nil {
    log.Fatal(err)
  }

  data, err := ioutil.ReadAll(resp.Body)
  resp.Body.Close()
  if err != nil {
    log.Fatal(err)
  }

  var igm InstagramMedia
  err = json.Unmarshal(data, &igm)
  if err != nil {
    log.Fatal(err)
  }

  author = igm.AuthorName

  client := instagram.NewClient(nil)
  client.ClientID = os.Getenv("CLIENT_ID")
  client.ClientSecret = os.Getenv("CLIENT_SECRET")
  client.AccessToken = tok

  media, err := client.Media.Get(igm.MediaId)
  if err != nil {
    log.Fatal(err)
  }

  videos := media.Videos
  if videos != nil {
    standard := videos.StandardResolution
    if standard == nil {
      log.Fatal("No valid media")
    }

    source = standard.URL
    if source == "" {
      log.Fatal("No valid media")
    }
  } else {
    images := media.Images
    standard := images.StandardResolution
    if standard == nil {
      log.Fatal("No valid media")
    }

    source = standard.URL
    if source == "" {
      log.Fatal("No valid media")
    }
  }

  return author, source
}

func get_ig_filename(path string) string {
  re := regexp.MustCompile("\\w+\\.(mp4|jpg|png)")
  filename := re.FindString(path)

  if os.Getenv("DEBUG") == "true" {
    fmt.Printf("[debug] get_ig_filename: path: %s, filename: %s\n", path, filename)
  }

  return filename
}

func save_source_to_dbox(author string, source string, client *dropbox.Dropbox) {
  filename := get_ig_filename(source)
  local := fmt.Sprintf("/tmp/%s", filename)

  // // https://github.com/thbar/golang-playground/blob/master/download-files.go
  output, err := os.Create(local)
  if err != nil {
    log.Fatal(fmt.Sprintf("Unable to create file: %s\n%s", filename, err))
  }
  defer output.Close()
  defer os.Remove(local)

  response, err := http.Get(source)
  if err != nil {
    log.Fatal(fmt.Sprintf("Unable to download file: %s\n%s", source, err))
  }
  defer response.Body.Close()

  _, err = io.Copy(output, response.Body)
  if err != nil {
    log.Fatal(fmt.Sprintf("Unable to download file: %s\n%s", source, err))
  }

  dst := fmt.Sprintf("%s/%s", author, filename)
  if _, err = client.UploadFile(local, dst, true, ""); err != nil {
    log.Fatal(fmt.Sprintf("Unable to upload file to dest: %s\n%s", dst, err))
  }
}

func main() {
  ig_access_token, err := ioutil.ReadFile("~/tmp/ig_access_token")

  var access_token string
  access_token = string(ig_access_token)

  if access_token == "" {
    token := get_ig_access_token()
    access_token = token.AccessToken
    err = ioutil.WriteFile("/tmp/ig_access_token", []byte(access_token), 0644)
    if err != nil {
      log.Fatal(err)
    }
  }

  db := dropbox.NewDropbox()
  db.SetAppInfo(os.Getenv("DB_CLIENT_ID"), os.Getenv("DB_CLIENT_SECRET"))
  db.SetAccessToken(os.Getenv("DB_ACCESS_TOKEN"))

  db.DownloadToFile("liked.txt", "/tmp/ig_liked", "")
  defer os.Remove("/tmp/ig_liked")
  defer db.Delete("liked.txt")

  liked, err := ioutil.ReadFile("/tmp/ig_liked")
  if err != nil {
    log.Fatal(err)
  }

  media_urls := strings.Split(string(liked), "\n")
  for _, url := range media_urls {
    if url != "" {
      fmt.Printf("Saving URL: %s\n", url)

      author, source := get_ig_source(url, access_token)
      save_source_to_dbox(author, source, db)
    }
  }
}
