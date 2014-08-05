package repository

import (
  "fmt"
  "github.com/chischaschos/twitnalytics/twitter"
)

func ExampleTweetsByUser() {
  Clear()

  tweets := []twitter.Tweet{
    twitter.Tweet{"text super cool", 1},
    twitter.Tweet{"not that cool", 5},
  }

  CreateTweets("chischaschos", tweets)

  TweetsByUser("chischaschos", func(tweet *twitter.Tweet) {
    fmt.Printf("%v\n", tweet)
  })

  // Output:
  // &{text super cool 1}
  // &{not that cool 5}
}
