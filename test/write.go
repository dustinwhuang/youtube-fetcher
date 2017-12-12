package main

import (
  "os"
  "fmt"
  "strconv"
  "bufio"
)

func main() {
  fmt.Println("Running fetcher...")

  f, err := os.Create("writes.csv")
  if err != nil {
    panic(err)
  }
  defer f.Close()

  w := bufio.NewWriter(f)
  for i := 0; i < 10000000; i++ {
    _, err = w.WriteString(strconv.Itoa(i) + "\n")
    if err != nil {
      panic(err)
    }
  }
  w.Flush()
}