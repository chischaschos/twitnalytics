package twitter

import (
  "fmt"
  "testing"
  "net/http"
  "net/http/httptest"
)

func TestFailedAuthenticate(t *testing.T) {
  testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, `{"errors":[{"code":99,"label":"authenticity_token_error","message":"Unable to verify your credentials"}]}`)
  }))

  twitter := Twitter{consumerKey: "ASD", consumerSecret: "ASD", endpoints: Endpoints{testServer.URL, ""}}
  authError := twitter.authenticate()

  if authError == nil {
    t.Error("An error was expected")
  } else if authError.Error() != "Unable to verify your credentials" {
    t.Error(authError.Error())
  }

  if twitter.accessToken != "" {
    t.Error("A token should not be returned")
  }

}

func TestSuccesfulAuthenticate(t *testing.T) {
  testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, `{"access_token":"i'm valid", "token_type":"bearer"}`)
  }))

  twitter := Twitter{consumerKey: "ASD", consumerSecret: "ASD", endpoints: Endpoints{testServer.URL, ""}}
  authError := twitter.authenticate()

  if authError != nil {
    t.Error("No error was expected, but:", authError.Error())
  }

  if twitter.accessToken != "i'm valid" {
    t.Error("A token should be returned", twitter.accessToken)
  }

}

func TestSuccessfulPullTweetsOf(t *testing.T) {
  authServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, `{"access_token":"i'm valid", "token_type":"bearer"}`)
  }))

  tweetsServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, `[{"text":"t1"}, {"text":"t3"}]`)
  }))

  twitter := Twitter{consumerKey: "ASD", consumerSecret: "ASD", endpoints: Endpoints{authServer.URL, tweetsServer.URL}}
  tweets, error := twitter.PullTweetsOf("chischaschos")

  if error != nil {
    t.Fatal("Tweets couldn't be pulled", error.Error())
  }

  if len(tweets) != 2 {
    t.Fatal("Wrong number of tweets retrieved")
  }

}
