package main

import (
  "os"
  "bufio"
  "net/http"
  "encoding/json"
  "io/ioutil"
  "fmt"
  "strings"
  "strconv"
  "time"
  "sync"
)

type Response struct {
  NextPageToken   string
  Items           []struct {
    Snippet       struct {
      ChannelId   string
    }
  }
}

func main() {
  if len(os.Args) < 6 {
    fmt.Println("Usage: ./youtube-fetcher <snippet_output_file> <contentDetails_output_file> <statistics_output_file> <topicDetails_output_file> <query_file>")
    os.Exit(2)
  }
  fmt.Println("Running fetcher...")
  start := time.Now()
  vids := 0

  f, err := os.OpenFile("tokens.csv", os.O_APPEND | os.O_WRONLY, 0622)
  if err != nil {
    panic(err)
  }
  defer f.Close()

  f1, err := os.OpenFile(os.Args[1], os.O_APPEND | os.O_WRONLY, 0622)
  if err != nil {
    panic(err)
  }
  defer f1.Close()

  f2, err := os.OpenFile(os.Args[2], os.O_APPEND | os.O_WRONLY, 0622)
  if err != nil {
    panic(err)
  }
  defer f2.Close()

  f3, err := os.OpenFile(os.Args[3], os.O_APPEND | os.O_WRONLY, 0622)
  if err != nil {
    panic(err)
  }
  defer f3.Close()

  f4, err := os.OpenFile(os.Args[4], os.O_APPEND | os.O_WRONLY, 0622)
  if err != nil {
    panic(err)
  }
  defer f4.Close()

  f5, err := os.OpenFile(os.Args[5], os.O_RDONLY, 0622)
  if err != nil {
    panic(err)
  }
  defer f5.Close()

  l, err := bufio.NewReader(f5).ReadString('\n')
  if err != nil {
    panic(err)
  }

  topics := strings.Split(l, ", ")

  client := &http.Client{}

  req, err := http.NewRequest("GET", "https://www.googleapis.com/youtube/v3/search", nil)
  if err != nil {
    panic(err)
  }


  q := req.URL.Query()
  q.Add("q", topics[0])
  topics = topics[1:]
  q.Add("part", "snippet")
  q.Add("chart", "mostPopular")
  q.Add("maxResults", "50")
  q.Add("key", os.Getenv("YOUTUBE_KEY"))
  q.Add("pageToken", "")

  mu := &sync.Mutex{}
  ch := make(chan string)
  go func() {
    for {
      id := <-ch
      go func(v string) {
        vids = vids + videos(v, mu, f1, f2, f3, f4)
      }(id)
    }
  }()

  c := make(map[string]bool)
  for ; vids < 10100100; {
    req.URL.RawQuery = q.Encode()
    resp, err := client.Do(req)
    if err != nil {
      panic(err)
    }

    respData, err := ioutil.ReadAll(resp.Body)
    respObject := Response{}
    json.Unmarshal(respData, &respObject)


    ids := []string{}
    for _, item := range respObject.Items {
      if _, ok := c[item.Snippet.ChannelId]; !ok {
          ids = append(ids, item.Snippet.ChannelId)
          c[item.Snippet.ChannelId] = true
        } 
    }

    ids = channels(strings.Join(ids, ","))

    for _, id := range ids {
      go func(p string) {
        for _, v := range playlists(p) {
          ch <- v
        }
      }(id)
    }

    resp.Body.Close()

    if respObject.NextPageToken == "" {
      if (len(topics) > 0) {
          q.Set("q", topics[0])
          f.WriteString(topics[0] + "\n")
          topics = topics[1:]
      } else {
        break
      }
    }
    q.Set("pageToken", respObject.NextPageToken)
    f.WriteString(respObject.NextPageToken + "\n")

    fmt.Printf("\rFetched " + strconv.Itoa(vids) + " videos in " + time.Now().Sub(start).String())
  }

  fmt.Println("\rFetched " + strconv.Itoa(vids) + " videos in " + time.Now().Sub(start).String())
}