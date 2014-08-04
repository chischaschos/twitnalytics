package data

import (
  "github.com/chischaschos/twitnalytics/twitter"
  "fmt"
)

func ExampleSyncTweets() {
  setupDB()
  clearDB()
  defer DB.Close()

  insertTestTweets()

  stmt, countError := DB.Query("SELECT count(*) from tweets")
  var count int

  if countError != nil {
    fmt.Println(countError)
  }

  stmt.Next()
  stmt.Scan(&count)
  fmt.Printf("%d tweets inserted\n", count)

  if count != 5 {
    fmt.Println(count)
  }

  rows, selectError := DB.Query("SELECT tweet_id, term, count FROM tweet_terms WHERE username = 'chischaschos'");

  if selectError != nil {
    fmt.Println(selectError)
  }

  for rows.Next() {
    var (
      tweet_id int64
      term string
      count int
    )
    rowError := rows.Scan(&tweet_id, &term, &count)

    if rowError != nil {
      fmt.Println(rowError)
    }

    fmt.Println(tweet_id, term, count)
  }

  // Output:
  // 5 tweets inserted
  // 1 text 4
  // 1 super 2
  // 1 cool 2
  // 1 bad 1
  // 1 is 2
  // 1 here 2
  // 1 not 2
  // 1 that 1

}

func clearDB() {
  _, execErr1 := DB.Exec("DELETE FROM tweets")
  _, execErr2 := DB.Exec("DELETE FROM tweet_terms")

  if execErr1 != nil {
    panic(execErr1)
  }

  if execErr2 != nil {
    panic(execErr2)
  }
}

func insertTestTweets() {
  tweets := []twitter.Tweet{
    twitter.Tweet{"text super cool", 1},
    twitter.Tweet{"text super bad", 2},
    twitter.Tweet{"text is here", 3},
    twitter.Tweet{"text is not here", 4},
    twitter.Tweet{"not that cool", 5},
  }

  SyncTweets("chischaschos", tweets)
}
