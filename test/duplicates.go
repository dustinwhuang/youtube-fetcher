package main

import (
  "os"
  "fmt"
  "bufio"
  "strings"
)

func main() {
  fmt.Println("Running fetcher...")

  f, err := os.Open("data.csv")
  if err != nil {
    panic(err)
  }
  defer f.Close()

  f1, err := os.Create("channels.csv")
  if err != nil {
    panic(err)
  }
  defer f1.Close()

  c := map[string]bool{}

  r := bufio.NewReader(f)
  w := bufio.NewWriter(f1)
  for i := 0; i < 164; i++ {
    l, err := r.ReadBytes('\n')
    if err != nil {
      panic(err)
    }

    id := strings.Split(string(l), ",")[2]

    if _, ok := c[id]; !ok {
      c[id] = true

      _, err = w.WriteString(id + ",")
      if err != nil {
        panic(err)
      }
    }
  }
  w.Flush()
}