package main

import (
	"encoding/json"
	"net/http"
	"os"
	"fmt"

	"github.com/garyburd/go-oauth/oauth"
	"github.com/pkg/errors"
)

const (
	refreshTokenURL  = "https://api.twitter.com/oauth/request_token"
	authorizationURL = "https://api.twitter.com/oauth/authenticate"
	accessTokenURL   = "https://api.twitter.com/oauth/access_token"
	accountURL       = "https://api.twitter.com/1.1/account/verify_credentials.json"

	callbackURL = "http://127.0.0.1:8080/login/twitter/auth/callback"
)

var (
	consumerKey    string
	consumerSecret string
)

func init() {
	consumerKey = os.Getenv("TEST_KEY")
	consumerSecret = os.Getenv("TEST_SECRET")
	if consumerSecret == "" {
		fmt.Println("env is not set")
		os.Exit(1)
	}
}

func NewTWClient() *oauth.Client {
	oc := &oauth.Client{
		TemporaryCredentialRequestURI: refreshTokenURL,
		ResourceOwnerAuthorizationURI: authorizationURL,
		TokenRequestURI:               accessTokenURL,
		Credentials: oauth.Credentials{
			Token:  consumerKey,
			Secret: consumerSecret,
		},
	}

	return oc
}

func GetAccessToken(rt *oauth.Credentials, oauthVerifier string) (int, *oauth.Credentials, error) {
	oc := NewTWClient()
	at, _, err := oc.RequestToken(nil, rt, oauthVerifier)
	if err != nil {
		err := errors.Wrap(err, "Failed to get access token.")
		return http.StatusBadRequest, nil, err
	}

	return http.StatusOK, at, nil
}

func GetMe(at *oauth.Credentials, user interface{}) (int, error) {
	oc := NewTWClient()
	resp, err := oc.Get(nil, at, accountURL, nil)
	if err != nil {
		err = errors.Wrap(err, "Failed to send twitter request.")
		return http.StatusInternalServerError, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 500 {
		err = errors.New("Twitter is unavailable")
		return http.StatusInternalServerError, err
	}

	if resp.StatusCode >= 400 {
		err = errors.New("Twitter request is invalid")
		return http.StatusBadRequest, err
	}

	err = json.NewDecoder(resp.Body).Decode(user)
	if err != nil {
		err = errors.Wrap(err, "Failed to decode user account response.")
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil

}
