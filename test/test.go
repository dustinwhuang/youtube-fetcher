package main

import (
  "net/http"
  "io/ioutil"
  "fmt"
)

func main() {
  fmt.Println("Running fetcher...")

  client := &http.Client{}

  req, err := http.NewRequest("GET", "https://www.googleapis.com/youtube/v3/videos", nil)
  if err != nil {
    panic(err)
  }

  q := req.URL.Query()
  q.Add("part", "snippet, statistics, topicDetails")
  q.Add("id", "2Vv-BfVoq4g")
  q.Add("maxResults", "50")
  q.Add("key", "AIzaSyAhsPBoJk-7lzrW7E5fsE5IzEntSuzgEdc")
  q.Add("pageToken", "")

  req.URL.RawQuery = q.Encode()
  resp, err := client.Do(req)
  if err != nil {
    panic(err)
  }

  respData, err := ioutil.ReadAll(resp.Body)
  fmt.Println(string(respData))


  resp.Body.Close()

}