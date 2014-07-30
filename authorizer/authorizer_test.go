package authorizer

import (
  "fmt"
  "testing"
  "net/http"
  "net/http/httptest"
)

func TestFailedAuthorize(t *testing.T) {
  testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, `{"errors":[{"code":99,"label":"authenticity_token_error","message":"Unable to verify your credentials"}]}`)
  }))

  authorizer := Authorization{ConsumerKey: "ASD", ConsumerSecret: "ASD", AccessToken: "", Tweets: nil, Endpoints: Endpoints{testServer.URL, ""}}
  authError := authorizer.Do()

  if authError == nil {
    t.Error("An error was expected")
  } else if authError.Error() != "Unable to verify your credentials" {
    t.Error(authError.Error())
  }

  if authorizer.AccessToken != "" {
    t.Error("A token should not be returned")
  }

}

func TestSuccesfulAuthorize(t *testing.T) {
  testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, `{"access_token":"i'm valid", "token_type":"bearer"}`)
  }))

  authorizer := Authorization{ConsumerKey: "ASD", ConsumerSecret: "ASD", AccessToken: "", Tweets: nil, Endpoints: Endpoints{testServer.URL, ""}}
  authError := authorizer.Do()

  if authError != nil {
    t.Error("No error was expected, but:", authError.Error())
  }

  if authorizer.AccessToken != "i'm valid" {
    t.Error("A token should be returned", authorizer.AccessToken)
  }

}

func TestPullTweets(t *testing.T) {
  authServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, `{"access_token":"i'm valid", "token_type":"bearer"}`)
  }))

  tweetsServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, `[{"text":"t1"}, {"text":"t3"}]`)
  }))

  authorizer := Authorization{ConsumerKey: "ASD", ConsumerSecret: "ASD", AccessToken: "", Tweets: nil, Endpoints: Endpoints{authServer.URL, tweetsServer.URL}}
  authError := authorizer.Do()

  if authError != nil {
    t.Fatal("Can't continue without an access token")
  }

  pullError := authorizer.PullTweets("chischaschos")

  if pullError != nil {
    t.Fatal("Tweets couldn't be pulled", pullError.Error())
  }
}
