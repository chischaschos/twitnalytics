package repository

import (
  "fmt"
  "database/sql"
  _ "code.google.com/p/go-sqlite/go1/sqlite3"
  "github.com/chischaschos/twitnalytics/twitter"
  "os"
  "os/user"
  "path/filepath"
)

func DB() *sql.DB {
  db, connErr := sql.Open("sqlite3", dbName())

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

  return db
}

func dbName() string {
  var dbName string

  if os.Getenv("TWITNALYTICS-ENV") == "test" {
    dbName = "../twitnalytics-db-test"

  } else {
    currentUser, userError := user.Current()

    if userError != nil {
      panic(userError)
    }

    dbName =  currentUser.HomeDir + "/.twitnalytics-db"
  }

  path, filepathError := filepath.Abs(dbName)

  if filepathError != nil {
    panic(filepathError)
  }

  return path
}

func TweetsByUser(username string, callback func(*twitter.Tweet))  {
  db := DB()

  rows, selectError := db.Query("SELECT id, text FROM tweets WHERE username = ?", username);

  if selectError != nil {
    panic(selectError)
  }

  for rows.Next() {
    var tweet twitter.Tweet

    rowError := rows.Scan(&tweet.Id, &tweet.Text)

    if rowError != nil {
      fmt.Printf("%#v\n", rowError)
    } else {
      callback(&tweet)
    }
  }
}

func TermsByUser(username string, callback func(*twitter.TermDoc)) {
  db := DB()
  rows, selectError := db.Query("SELECT tweet_id, term, count FROM tweet_terms WHERE username = ?", username);

  if selectError != nil {
    panic(selectError)
  }

  for rows.Next() {
    var termDoc twitter.TermDoc
    rowError := rows.Scan(&termDoc.TweetId, &termDoc.Term, &termDoc.Count)

    if rowError != nil {
      fmt.Println(rowError)
    } else {
      callback(&termDoc)
    }
  }
}

func CreateTweets(username string, tweets []twitter.Tweet) {
  db := DB()
  defer db.Close()

  stmt, prepareError := db.Prepare("INSERT INTO tweets VALUES(?, ?, ?)")

  if prepareError != nil {
    panic(prepareError)
  }

  for _, tweet := range tweets {
    _, execError := stmt.Exec(tweet.Id, tweet.Text, username)

    if execError != nil {
      fmt.Println(execError)
    }

  }
}

func SaveTerms(username string, termsDictionary map[string]*twitter.TermDoc) {
  db := DB()
  defer db.Close()

  insertStmt, stmtError := db.Prepare("INSERT INTO tweet_terms (username, tweet_id, term, count) VALUES (?, ?, ?, ?)")

  if stmtError != nil {
    panic(stmtError)
  }

  for _, termDoc := range termsDictionary {
    _, ie := insertStmt.Exec(username, termDoc.TweetId, termDoc.Term, termDoc.Count)

    if ie != nil {
      fmt.Println(ie)
    }
  }
}

func Clear() {
  db := DB()
  _, execErr1 := db.Exec("DELETE FROM tweets")
  _, execErr2 := db.Exec("DELETE FROM tweet_terms")

  if execErr1 != nil {
    panic(execErr1)
  }

  if execErr2 != nil {
    panic(execErr2)
  }
}

func TotalTweets() int {
  db := DB()
  stmt, countError := db.Query("SELECT count(*) from tweets")

  if countError != nil {
    panic(countError)
  }

  var count int

  stmt.Next()
  stmt.Scan(&count)

  return count
}
