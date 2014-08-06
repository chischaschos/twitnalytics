package main

import (
  "flag"
  "fmt"
  tw "github.com/chischaschos/twitnalytics/twitter"
  "github.com/chischaschos/twitnalytics/repository"
  "github.com/chischaschos/twitnalytics/data"
  "os"
  "os/user"
  "path/filepath"
)

var username string

func init() {
  flag.StringVar(&username, "u", "the_user_name", "the user name whose timeline we are gonna play with")

  repository.Settings().StorePath = settingsPath()
}

func settingsPath() string {
  currentUser, userError := user.Current()

  if userError != nil {
    panic(userError)
  }

  settingsFile := currentUser.HomeDir + "/.twitnalytics.json"
  settingsPath, filepathError := filepath.Abs(settingsFile)

  if filepathError != nil {
    panic(filepathError)
  }

  return settingsPath
}

func authValues() (string, string) {
  consumerKey := repository.Settings().Get("consumer-key")

  if consumerKey == "" {
    consumerKey := os.Getenv("CONSUMER_KEY")

    if consumerKey == "" {
      panic("Twitter CONSUMER_KEY not defined, please export it")
    } else {
      repository.Settings().Set("consumer-key", consumerKey)
    }
  }

  consumerSecret := repository.Settings().Get("consumer-secret")

  if consumerSecret == "" {
    consumerSecret := os.Getenv("CONSUMER_SECRET")

    if consumerSecret == "" {
      panic("Twitter CONSUMER_SECRET not defined, please export it")
    } else {
      repository.Settings().Set("consumer-secret", consumerSecret)
    }
  }

  return consumerKey, consumerSecret
}

func main() {
  flag.Parse()

  if username == "the_user_name" {
    flag.PrintDefaults()
  } else {

    consumerKey, consumerSecret := authValues()

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
