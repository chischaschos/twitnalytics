package main

import (
  "flag"
  "fmt"
  tw "github.com/chischaschos/twitnalytics/twitter"
  "github.com/chischaschos/twitnalytics/repository"
  "github.com/chischaschos/twitnalytics/data"
  "os"
)

var username string

func init() {
  flag.StringVar(&username, "u", "the_user_name", "the user name whose timeline we are gonna play with")
}

func main() {
  flag.Parse()

  if username == "the_user_name" {
    flag.PrintDefaults()
  } else {

    consumerKey := os.Getenv("CONSUMER_KEY")
    consumerSecret := os.Getenv("CONSUMER_SECRET")

    twitter := tw.New(consumerKey, consumerSecret)
    tweets, pullError := twitter.PullTweetsOf(username)

    fmt.Println(pullError)

    for _, tweet := range tweets {
      fmt.Println(tweet)
    }

    if len(tweets) > 0 {
      data.SyncTweets(username, tweets)
    }

    repository.TermsByUser(username, func(termDoc *tw.TermDoc) {
      fmt.Printf("%v\n", termDoc)
    })
  }
}
