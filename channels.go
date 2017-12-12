package main

import (
  "os"
  "net/http"
  "encoding/json"
  "io/ioutil"
)

type ChannelResponse struct {
  NextPageToken         string
  Items                 []struct {
    ContentDetails      struct {
      RelatedPlaylists  struct {
        Uploads         string
      }
    }
  }
}

func channels(ids string) []string {
  client := &http.Client{}

  req, err := http.NewRequest("GET", "https://www.googleapis.com/youtube/v3/channels", nil)
  if err != nil {
    panic(err)
  }

  q := req.URL.Query()
  q.Add("part", "contentDetails")
  q.Add("id", ids)
  q.Add("maxResults", "50")
  q.Add("key", os.Getenv("YOUTUBE_KEY"))
  req.URL.RawQuery = q.Encode()

  resp, err := client.Do(req)
  if err != nil {
    panic(err)
  }

  respData, err := ioutil.ReadAll(resp.Body)
  respObject := ChannelResponse{}
  json.Unmarshal(respData, &respObject)

  p := []string{}

  for _, item := range respObject.Items {
    p = append(p, item.ContentDetails.RelatedPlaylists.Uploads)
  }

  resp.Body.Close()

  return p
}