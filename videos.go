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
  Id      string      `json:"id"`
  Items   []Items     `json:"items"`
}

type Items struct {
  Id              string          `json:"id"`
  Snippet         Snippet         `json:"snippet"`
  ContentDetails  ContentDetails  `json:"contentDetails"`
  Statistics      Statistics      `json:"statistics"`
  TopicDetails    TopicDetails    `json:"topicDetails"`
}

type Snippet struct {
  Id            string        `json:"_id"`
  PublishedAt   string        `json:"publishedAt"`
  ChannelId     string        `json:"channelId"`
  Title         string        `json:"title"`
  Description   string        `json:"description"`
  Thumbnails    struct {
    Default     struct {
      Url       string        `json:"url"`
    } `json:"default"`
  } `json:"thumbnails"`
  Tags          []string      `json:"tags"`
  CategoryId    string        `json:"categoryId"`
}

type ContentDetails struct {
  Id            string        `json:"_id"`
  Duration      string        `json:"duration"`
}

type Statistics  struct {
  Id             string       `json:"_id"`
  ViewCount      string       `json:"viewCount"`
}

type TopicDetails struct {
  Id                string    `json:"_id"`
  RelevantTopicIds  []string  `json:"relevantTopicIds"`
  TopicCategories   []string  `json:"topicCategories"`
}

func videos(id string, mu *sync.Mutex, f1 *os.File, f2 *os.File, f3 *os.File, f4 *os.File) int {
  client := &http.Client{}

  req, err := http.NewRequest("GET", "https://www.googleapis.com/youtube/v3/videos", nil)
  if err != nil {
    panic(err)
  }

  q := req.URL.Query()
  q.Add("part", "snippet, contentDetails, statistics, topicDetails")
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

  i := 0
  item := Items{}
  mu.Lock()
  w1 := bufio.NewWriter(f1)
  w2 := bufio.NewWriter(f2)
  w3 := bufio.NewWriter(f3)
  w4 := bufio.NewWriter(f4)
  for i, item = range respObject.Items {
    item.Snippet.Id = item.Id
    str, err := json.Marshal(item.Snippet)
    _, err = w1.WriteString(string(str) + "\n")
    if err != nil {
      panic(err)
    }

    item.ContentDetails.Id = item.Id
    str, err = json.Marshal(item.ContentDetails)
    _, err = w2.WriteString(string(str) + "\n")
    if err != nil {
      panic(err)
    }

    item.Statistics.Id = item.Id
    str, err = json.Marshal(item.Statistics)
    _, err = w3.WriteString(string(str) + "\n")
    if err != nil {
      panic(err)
    }

    item.TopicDetails.Id = item.Id
    str, err = json.Marshal(item.TopicDetails)
    _, err = w4.WriteString(string(str) + "\n")
    if err != nil {
      panic(err)
    }
  }
  w1.Flush()
  w2.Flush()
  w3.Flush()
  w4.Flush()
  mu.Unlock()

  resp.Body.Close()

  return i
}