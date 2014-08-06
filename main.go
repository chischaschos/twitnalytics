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
  "strings"
)

var usernames string

func init() {
  flag.StringVar(&usernames, "u", "user1,user2,user3", "the user names whose timelines we are gonna play with")

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

  if usernames == "user1,user2,user3" {
    flag.PrintDefaults()
  } else {

    consumerKey, consumerSecret := authValues()
    twitter := tw.New(consumerKey, consumerSecret)

    tweetsChannel := make(chan map[string][]tw.Tweet)
    splitUsernames := strings.Split(usernames, ",")

    for _, username := range splitUsernames {
      go twitter.PullTweetsOf(username, tweetsChannel)
    }

    usersWithTweets := joinTweets(tweetsChannel, len(splitUsernames))

    for _, userWithTweets := range usersWithTweets {
      for username, tweets := range userWithTweets {

        if len(tweets) > 0 {
          data.SyncTweets(username, tweets)
        }

        repository.TermsByUser(username, func(termDoc *tw.TermDoc) {
          fmt.Printf("%v\n", termDoc)
        })
      }
    }
  }
}

func joinTweets(tweetsChannel chan map[string][]tw.Tweet, usersCount int) []map[string][]tw.Tweet {
  var usersWithTweets []map[string][]tw.Tweet

  for {
    select {
    case message := <-tweetsChannel:
      usersWithTweets = append(usersWithTweets, message)

      if len(usersWithTweets) == usersCount {
        return usersWithTweets
      }
    }
  }
}
