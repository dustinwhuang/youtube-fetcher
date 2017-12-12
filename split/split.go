package main

import (
  "os"
  "io"
  "bufio"
  "encoding/json"
  "time"
  "fmt"
  "strconv"
)

type Video struct {
  Id              string          `json:"id"`
  Snippet         Snippet         `json:"snippet"`
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

type Statistics  struct {
  Id             string        `json:"_id"`
  ViewCount      string       `json:"viewCount"`
}

type TopicDetails struct {
  Id                string    `json:"_id"`
  RelevantTopicIds  []string  `json:"relevantTopicIds"`
  TopicCategories   []string  `json:"topicCategories"`
}

func main() {
  fmt.Println("Running split...")
  start := time.Now()
  rows := 0

  f, err := os.OpenFile(os.Args[1], os.O_RDONLY, 0622)
  if err != nil {
    panic(err)
  }
  defer f.Close()

  f1, err := os.OpenFile("../data/snippet.json", os.O_APPEND | os.O_WRONLY, 0622)
  if err != nil {
    panic(err)
  }
  defer f1.Close()

  f2, err := os.OpenFile("../data/statistics.json", os.O_APPEND | os.O_WRONLY, 0622)
  if err != nil {
    panic(err)
  }
  defer f2.Close()

  f3, err := os.OpenFile("../data/topicDetails.json", os.O_APPEND | os.O_WRONLY, 0622)
  if err != nil {
    panic(err)
  }
  defer f3.Close()

  r := bufio.NewReader(f)
  w1 := bufio.NewWriter(f1)
  w2 := bufio.NewWriter(f2)
  w3 := bufio.NewWriter(f3)
  respObject := Video{}
  for {
    l, _, err := r.ReadLine()
    if err != nil {
      if err == io.EOF {
        break
      }
      panic(err)
    }

    json.Unmarshal(l, &respObject)

    respObject.Snippet.Id = respObject.Id
    str, err := json.Marshal(respObject.Snippet)
    _, err = w1.WriteString(string(str) + "\n")
    if err != nil {
      panic(err)
    }

    respObject.Statistics.Id = respObject.Id
    str, err = json.Marshal(respObject.Statistics)
    _, err = w2.WriteString(string(str) + "\n")
    if err != nil {
      panic(err)
    }

    respObject.TopicDetails.Id = respObject.Id
    str, err = json.Marshal(respObject.TopicDetails)
    _, err = w3.WriteString(string(str) + "\n")
    if err != nil {
      panic(err)
    }

    rows++
    if (rows % 10000 == 0) {
      fmt.Printf("\rSplit " + strconv.Itoa(rows) + " rows in " + time.Now().Sub(start).String())
    }
  }
  w1.Flush()
  w2.Flush()
  w3.Flush()

  fmt.Println("\rSplit " + strconv.Itoa(rows) + " rows in " + time.Now().Sub(start).String())
}
