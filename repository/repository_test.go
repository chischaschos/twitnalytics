package repository

import (
  "fmt"
  "github.com/chischaschos/twitnalytics/twitter"
  "os"
)

func ExampleTweetsByUser() {
  os.Setenv("TWITNALYTICS-ENV", "test")

  Clear()

  tweets := []twitter.Tweet{
    twitter.Tweet{"text super cool", 496697195907665921},
    twitter.Tweet{"not that cool", 4966971959076659215},
  }

  CreateTweets("chischaschos", tweets)

  TweetsByUser("chischaschos", func(tweet *twitter.Tweet) {
    fmt.Printf("%v\n", tweet)
  })

  // Output:
  // &{text super cool 496697195907665921}
  // &{not that cool 4966971959076659215}
}
