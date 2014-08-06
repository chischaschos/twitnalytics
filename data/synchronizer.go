package data

import (
  "github.com/chischaschos/twitnalytics/repository"
  "github.com/chischaschos/twitnalytics/twitter"
  "strings"
)

func SyncTweets(username string, tweets []twitter.Tweet) {
  repository.CreateTweets(username, tweets)

  quit := make(chan int)

  // The idea here will be to extracts chunks of terms one chunk per goroutine
  go extractTermsForUser(username, quit)

  for {
    select {
    case <-quit:
      return
    }
  }
}

func extractTermsForUser(username string, quit chan<- int) {
  termsDictionary := map[string]*twitter.TermDoc{}

  repository.TweetsByUser(username, func(tweet *twitter.Tweet) {
    terms := strings.Split(tweet.Text, " ")

    for _, term := range terms {
      _, ok := termsDictionary[term]

      if ok {
        termsDictionary[term].Count++
      } else {
        termsDictionary[term] = &twitter.TermDoc{TweetId: tweet.Id, Term: term, Count: 1}
      }
    }
  })

  repository.SaveTerms(username, termsDictionary)

  quit <- 1
}
