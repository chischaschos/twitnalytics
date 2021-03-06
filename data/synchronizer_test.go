package data

import (
  "github.com/chischaschos/twitnalytics/repository"
  "github.com/chischaschos/twitnalytics/twitter"
  "fmt"
  "os"
)

func ExampleSyncTweets() {
  os.Setenv("TWITNALYTICS-ENV", "test")

  repository.Clear()

  tweets := []twitter.Tweet{
    twitter.Tweet{"text super cool", 1},
    twitter.Tweet{"text super bad", 2},
    twitter.Tweet{"text is here", 3},
    twitter.Tweet{"text is not here", 4},
    twitter.Tweet{"not that cool", 5},
  }

  SyncTweets("chischaschos", tweets)

  fmt.Printf("%d tweets inserted\n", repository.TotalTweets())

  repository.TermsByUser("chischaschos", func(termDoc *twitter.TermDoc) {
    fmt.Printf("%v\n", termDoc)
  })

  // Output:
  // 5 tweets inserted
  // &{1 text 4}
  // &{1 super 2}
  // &{1 cool 2}
  // &{2 bad 1}
  // &{3 is 2}
  // &{3 here 2}
  // &{4 not 2}
  // &{5 that 1}
}

func ExampleNormalizeText() {
  fmt.Println(normalizeText("ASD"))
  // Output:
  // asd
}
