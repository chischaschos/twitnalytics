package main

import (
  "flag"
  "fmt"
  "github.com/chischaschos/twitnalytics/twitter"
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

    twitter := twitter.New(consumerKey, consumerSecret)
    tweets, pullError := twitter.PullTweetsOf(username)

    fmt.Println(pullError)

    for _, tweet := range tweets {
      fmt.Println(tweet)
    }
  }
}
