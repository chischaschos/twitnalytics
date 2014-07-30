package main

import (
  "fmt"
  "github.com/chischaschos/twitnalytics/authorizer"
  "os"
)

func main() {
  fmt.Println("On")
  consumerKey := os.Getenv("CONSUMER_KEY")
  consumerSecret := os.Getenv("CONSUMER_SECRET")

  authorizer := authorizer.New(consumerKey, consumerSecret)
  authorizer.Do()
  pullError := authorizer.PullTweets("chischaschos")

  if pullError != nil {
    fmt.Printf("Problem!!!: %s\n", pullError.Error())
  }

  for _, tweet := range authorizer.Tweets {
    fmt.Println(tweet)
  }
  //tweetsWithSimilarity := authorization.CalculateSimilarity("chischaschos")
}
