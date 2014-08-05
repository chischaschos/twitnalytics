package twitter

import (
  "encoding/base64"
  "encoding/json"
  "errors"
  "bytes"
  "net/http"
)

const defaultAuthenticationEndpoint = "https://api.twitter.com/oauth2/token"
const defaultTweetsEndpoint = "https://api.twitter.com/1.1/statuses/user_timeline.json"
const authBody = "grant_type=client_credentials"
const contentType = "application/x-www-form-urlencoded;charset=UTF-8"

type Endpoints struct {
  authentication, tweets string
}

type Twitter struct {
  consumerKey, consumerSecret string
  accessToken string
  endpoints Endpoints
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
  Id int64
}

type TermDoc struct {
  TweetId int64
  Term string
  Count int
}

// PullTweetsOf the retrieves a parameter user tweets
func (twitter *Twitter) PullTweetsOf(user string) ([]Tweet, error) {
  authError := twitter.authenticate()

  if authError != nil {
    return nil, authError
  }

  request, requestCreationError := http.NewRequest("GET", twitter.endpoints.tweets + "?screen_name=" + user, nil)

  if requestCreationError != nil {
    return nil, requestCreationError
  }

  request.Header.Add("Authorization", "Bearer " + twitter.accessToken)

  client := new(http.Client)
  response, requestError := client.Do(request)

  if requestError != nil {
    return nil, requestError
  }

  var tweets []Tweet
  jsonDecoder := json.NewDecoder(response.Body)
  jsonDecoderError := jsonDecoder.Decode(&tweets)

  if jsonDecoderError != nil {
    return nil, jsonDecoderError
  }

  return tweets, nil
}
// New creates a new Twitter client
func New(consumerKey, consumerSecret string) *Twitter {
  return &Twitter{consumerKey: consumerKey, consumerSecret: consumerSecret, endpoints: Endpoints{defaultAuthenticationEndpoint, defaultTweetsEndpoint}}
}

// authenticate is to be used internally to retrieve a client_credentials,
// grant type before pulling a user tweets
func (twitter *Twitter) authenticate() (authError error) {
  response, requestError := twitter.doAuthenticationRequest()

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
    twitter.accessToken = accessToken
  }

  return responseError
}

// doAuthenticationRequest makes the actual authentication call
func (twitter *Twitter) doAuthenticationRequest() (response *http.Response, requestError error) {
  request, requestCreationError := http.NewRequest("POST", twitter.endpoints.authentication, bytes.NewBufferString(authBody))

  if requestCreationError != nil {
    return nil, requestCreationError
  }

  authorizationToken := twitter.consumerKey + ":" + twitter.consumerSecret
  encodedAuthorizationToken := base64.StdEncoding.EncodeToString([]byte(authorizationToken))
  request.Header.Add("Authorization", "Basic " + encodedAuthorizationToken)
  request.Header.Add("Content-Type", contentType)

  client := new(http.Client)
  return client.Do(request)
}

// extractMessage maps the error messages into one big error message
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


