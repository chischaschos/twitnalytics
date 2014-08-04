package data

import (
  "github.com/chischaschos/twitnalytics/twitter"
  "database/sql"
  _ "code.google.com/p/go-sqlite/go1/sqlite3"
  "fmt"
  "strings"
)

var DB *sql.DB

func init() {
  setupDB()
}

func SyncTweets(username string, tweets []twitter.Tweet) {
  insertTweets(username, tweets)

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

func setupDB() {
  db, connErr := sql.Open("sqlite3", "twitnalytics")

  if connErr != nil {
    panic(connErr)
  }

  _, execErr1 := db.Exec("CREATE TABLE IF NOT EXISTS tweets(id, text, username)")
  _, execErr2 := db.Exec("CREATE TABLE IF NOT EXISTS tweet_terms(username, tweet_id, term, count)")

  if execErr1 != nil {
    panic(execErr1)
  }

  if execErr2 != nil {
    panic(execErr2)
  }

  DB = db
}

func extractTermsForUser(username string, quit chan<- int) {
  setupDB()

  rows, selectError := DB.Query("SELECT id, text FROM tweets WHERE username = ?", username);

  if selectError != nil {
    panic(selectError)
  }

  termsDictionary := map[string]int{}

  for rows.Next() {
    var (
      id int64
      text string
    )

    rowError := rows.Scan(&id, &text)

    terms := strings.Split(text, " ")

    for _, term := range terms {

      if rowError == nil {
        _, ok := termsDictionary[term]

        if ok {
          termsDictionary[term]++
        } else {
          termsDictionary[term] = 1
        }
      } else {
        fmt.Printf("%#v\n", rowError)
      }
    }
  }

  saveTerms(username, termsDictionary)

  quit <- 1
}

func insertTweets(username string, tweets []twitter.Tweet) {
  stmt, prepareError := DB.Prepare("INSERT INTO tweets VALUES(?, ?, ?)")

  if prepareError != nil {
    fmt.Println(prepareError)
  }

  for _, tweet := range tweets {
    _, execError := stmt.Exec(tweet.Id, tweet.Text, username)

    if execError != nil {
      fmt.Println(execError)
    }

  }
}

func saveTerms(username string, termsDictionary map[string]int) {
  insertStmt, stmtError := DB.Prepare("INSERT INTO tweet_terms (username, tweet_id, term, count) VALUES (?, ?, ?, ?)")

  if stmtError != nil {
    fmt.Println(stmtError)
  }

  for key, value := range termsDictionary {
    _, ie := insertStmt.Exec(username, 1, key, value)

    if ie != nil {
      fmt.Println(ie)
    }
  }
}
