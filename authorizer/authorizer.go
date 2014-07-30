package authorizer

import (
  "encoding/base64"
  "encoding/json"
  "errors"
  "bytes"
  "net/http"
)

const defaultAuthenticationEndpoint = "https://api.twitter.com/oauth2/token"
const defaultTweetsEndpoint = "https://api.twitter.com/1.1/statuses/user_timeline.json?screen_name="
const authBody = "grant_type=client_credentials"
const contentType = "application/x-www-form-urlencoded;charset=UTF-8"

type Endpoints struct {
  Authentication, Tweets string
}

type Authorization struct {
  ConsumerKey, ConsumerSecret, AccessToken string
  Tweets []Tweet
  Endpoints
}

type AuthorizationResponse struct {
  TokenType string `json:"token_type"`
  AccessToken string `json:"access_token"`
  Errors []struct {
    Code int
    Label string
    Message string
  }
}

type Tweet struct {
  Text string
}

func New(consumerKey, consumerSecret string) (authorization *Authorization) {
  return &Authorization{consumerKey, consumerSecret, "", nil,
    Endpoints{defaultAuthenticationEndpoint, defaultTweetsEndpoint}}
}

func (authorization *Authorization) Do() (authError error) {
  response, requestError := authorization.callEndpoint()

  if requestError != nil {
    return requestError
  }

  var authorizationResponse AuthorizationResponse
  jsonDecoder := json.NewDecoder(response.Body)
  jsonDecoderError := jsonDecoder.Decode(&authorizationResponse)

  if jsonDecoderError != nil {
    return jsonDecoderError
  }

  accessToken, responseError := extractMessage(authorizationResponse)

  if responseError == nil {
    authorization.AccessToken = accessToken
  }

  return responseError
}

func (authorization *Authorization) callEndpoint() (response *http.Response, requestError error) {
  request, requestCreationError := http.NewRequest("POST", authorization.Endpoints.Authentication, bytes.NewBufferString(authBody))

  if requestCreationError != nil {
    return nil, requestCreationError
  }

  authorizationToken := authorization.ConsumerKey + ":" + authorization.ConsumerSecret
  encodedAuthorizationToken := base64.StdEncoding.EncodeToString([]byte(authorizationToken))
  request.Header.Add("Authorization", "Basic " + encodedAuthorizationToken)
  request.Header.Add("Content-Type", contentType)

  client := new(http.Client)
  return client.Do(request)
}

func extractMessage(authorizationResponse AuthorizationResponse) (message string, responseError error) {

  if len(authorizationResponse.Errors) != 0 {
    var errorMessage string

    for _, authError := range authorizationResponse.Errors {
      errorMessage += authError.Message
    }

    return "", errors.New(errorMessage)

  } else {
    return authorizationResponse.AccessToken, nil
  }
}

func (authorization *Authorization) PullTweets(user string) error {
  request, requestCreationError := http.NewRequest("GET", authorization.Endpoints.Tweets + user, nil)

  if requestCreationError != nil {
    return requestCreationError
  }

  request.Header.Add("Authorization", "Bearer " + authorization.AccessToken)

  client := new(http.Client)
  response, requestError := client.Do(request)

  if requestError != nil {
    return requestError
  }

  var tweets []Tweet
  jsonDecoder := json.NewDecoder(response.Body)
  jsonDecoderError := jsonDecoder.Decode(&tweets)

  if jsonDecoderError != nil {
    return jsonDecoderError
  }

  authorization.Tweets = tweets

  return nil
}
