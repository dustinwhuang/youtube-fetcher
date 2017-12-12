package main

import (
  "os"
  "net/http"
  "encoding/json"
  "io/ioutil"
  "strings"
)

type PlaylistResponse struct {
  NextPageToken         string
  Items   []struct {
    ContentDetails    struct {
      VideoId         string
    }
  }
}

func playlists(id string) []string {
  client := &http.Client{}

  req, err := http.NewRequest("GET", "https://www.googleapis.com/youtube/v3/playlistItems", nil)
  if err != nil {
    panic(err)
  }

  q := req.URL.Query()
  q.Add("part", "contentDetails")
  q.Add("playlistId", id)
  q.Add("maxResults", "50")
  q.Add("key", os.Getenv("YOUTUBE_KEY"))
  q.Add("nextPageToken", "")

  l := []string{}
  for {
    req.URL.RawQuery = q.Encode()
    resp, err := client.Do(req)
    if err != nil {
      panic(err)
    }

    respData, err := ioutil.ReadAll(resp.Body)
    respObject := PlaylistResponse{}
    json.Unmarshal(respData, &respObject)

    v := []string{}

    for _, item := range respObject.Items {
      v = append(v, item.ContentDetails.VideoId)
    }

    resp.Body.Close()

    l = append(l, strings.Join(v, ","))

    if respObject.NextPageToken == "" {
      break;
    }
    q.Set("pageToken", respObject.NextPageToken)
  }

  return l
}