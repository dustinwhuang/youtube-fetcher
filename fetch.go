package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Response struct {
	NextPageToken string
	Items         []struct {
		Snippet struct {
			ChannelId string
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

	f := [6]*os.File{}
	var err error

	f[0], err = os.OpenFile("tokens.log", os.O_APPEND|os.O_WRONLY, 0622)
	if err != nil {
		panic(err)
	}
	defer f[0].Close()

	for i := 1; i < 5; i++ {
		if _, err := os.Stat(os.Args[i]); os.IsNotExist(err) {
			f[i], err = os.Create(os.Args[i])
		} else {
			f[i], err = os.OpenFile(os.Args[i], os.O_APPEND|os.O_WRONLY, 0622)
			if err != nil {
				panic(err)
			}
		}
		defer f[i].Close()
	}

	f[5], err = os.OpenFile(os.Args[5], os.O_RDONLY, 0622)
	if err != nil {
		panic(err)
	}
	defer f[5].Close()

	l, err := bufio.NewReader(f[5]).ReadString('\n')
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
				vids = vids + videos(v, mu, f[1], f[2], f[3], f[4])
			}(id)
		}
	}()

	cIds := make(map[string]bool)
	for vids < 1 {
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
			if _, ok := cIds[item.Snippet.ChannelId]; !ok {
				ids = append(ids, item.Snippet.ChannelId)
				cIds[item.Snippet.ChannelId] = true
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
			if len(topics) > 0 {
				q.Set("q", topics[0])
				f[0].WriteString(topics[0] + "\n")
				topics = topics[1:]
			} else {
				break
			}
		}
		q.Set("pageToken", respObject.NextPageToken)
		f[0].WriteString(respObject.NextPageToken + "\n")

		fmt.Printf("\rFetched " + strconv.Itoa(vids) + " videos in " + time.Now().Sub(start).String())
	}

	fmt.Println("\rFetched " + strconv.Itoa(vids) + " videos in " + time.Now().Sub(start).String())
}
