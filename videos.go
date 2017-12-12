package main

import (
  "os"
  "net/http"
  "encoding/json"
  "io/ioutil"
  "bufio"
  "sync"
)

type VideoResponse struct {
  Items   []Items     `json:"items"`
}

func videos(id string, mu *sync.Mutex, f *os.File) int {
  client := &http.Client{}

  req, err := http.NewRequest("GET", "https://www.googleapis.com/youtube/v3/videos", nil)
  if err != nil {
    panic(err)
  }

  q := req.URL.Query()
  q.Add("part", "snippet, statistics, topicDetails")
  q.Add("id", id)
  q.Add("maxResults", "50")
  q.Add("key", os.Getenv("YOUTUBE_KEY"))
  req.URL.RawQuery = q.Encode()

  resp, err := client.Do(req)
  if err != nil {
    panic(err)
  }

  respData, err := ioutil.ReadAll(resp.Body)
  respObject := VideoResponse{}
  json.Unmarshal(respData, &respObject)

  var i int
  var item Items
  mu.Lock()
  w := bufio.NewWriter(f)
  for i, item = range respObject.Items {
    str, err := json.Marshal(item)
    _, err = w.WriteString(string(str) + "\n")
    if err != nil {
      panic(err)
    }
  }
  w.Flush()
  mu.Unlock()

  resp.Body.Close()

  return i
}