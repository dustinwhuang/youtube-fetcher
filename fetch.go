package main

import (
  "os"
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
  NextPageToken string    `json:"nextPageToken"`
  Items         []Items   `json:"items"`
}

type Items struct {
  Id              string          `json:"id"`
  Snippet         Snippet         `json:"snippet"`
  ContentDetails  ContentDetails  `json:"contentDetails"`
  Statistics      Statistics      `json:"statistics"`
  TopicDetails    TopicDetails    `json:"topicDetails"`
}

type Snippet struct {
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
  Duration      string        `json:"duration"`
}

type Statistics  struct {
  ViewCount      string       `json:"viewCount"`
}

type TopicDetails struct {
  RelevantTopicIds  []string  `json:"relevantTopicIds"`
  TopicCategories   []string  `json:"topicCategories"`
}

func main() {
  if (len(os.Args) < 2) {
    fmt.Println("Usage: ./youtube-fetcher <output_file>")
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

  client := &http.Client{}

  req, err := http.NewRequest("GET", "https://www.googleapis.com/youtube/v3/search", nil)
  if err != nil {
    panic(err)
  }
  // old topics: "Autos & Vehicle", " Film & Animation", "Music", "Pets & Animals", "Sports", "Short Movies", "Gaming", "Videoblogging", "People & Blogs", "Comedy", "Entertainment", "News & Politics", "Howto & Style", "Education", "Nonprofits & Activism", "Movies", "Anime/Animation", "Action/Adventure", "Classics", "Comedy", "Documentary", "Drama", "Foreign", "Horror", "Sci-Fi/Fantasy", "Thriller", "Shorts", "Shows", "Trailers", "JavaScript", 
  topics := []string{"Magic", "Dance"}

  q := req.URL.Query()
  q.Add("q", "Pokemon")
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
        vids = vids + videos(v, mu, f1)
      }(id)
    }
  }()

  c := make(map[string]bool)

  for ; vids < 1; {
    req.URL.RawQuery = q.Encode()
    resp, err := client.Do(req)
    if err != nil {
      panic(err)
    }

    respData, err := ioutil.ReadAll(resp.Body)
    respObject := Response{}
    json.Unmarshal(respData, &respObject)


    for _, item := range respObject.Items {
      c[item.Snippet.ChannelId] = true
    }

    resp.Body.Close()

    if respObject.NextPageToken == "" {
      ids := []string{}
      for k := range c {
        ids = append(ids, k)
      }

      arr := [][]string{}
      for {
        arr = append(arr, ids[0:50])
        if (len(ids) > 50) {
          ids = ids[50:]
        } else {
          break
        }
      }
      
      for _, ids := range arr {
        ids = channels(strings.Join(ids, ","))

        for _, id := range ids {
          go func(p string) {
            for _, v := range playlists(p) {
              ch <- v
            }
          }(id)
        }
      }
      c = make(map[string]bool)

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