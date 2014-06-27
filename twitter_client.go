package main

import (
	"fmt"
	"github.com/mrjones/oauth"
	"log"
)

// TwitterClient type
type TwitterClient struct {
	consumer    *oauth.Consumer
	accessToken *oauth.AccessToken
}

// NewTwitterClient creates new Twitter client
func NewTwitterClient() *TwitterClient {
	return &TwitterClient{
		consumer: oauth.NewConsumer(
			"St8EoKgvUJEZ71o7GyvmSxP5w",
			"n7eS5ZDJELZLv466FYzpbMwSc9xmLkUvTkvy6AxWfvEjJZ4zWj",
			oauth.ServiceProvider{
				RequestTokenUrl:   "https://api.twitter.com/oauth/request_token",
				AuthorizeTokenUrl: "https://api.twitter.com/oauth/authorize",
				AccessTokenUrl:    "https://api.twitter.com/oauth/access_token",
			},
		),
	}
}

// PrepareAccessToken sets valid access token
func (c *TwitterClient) PrepareAccessToken() (err error) {
	conf, err := NewConfig("go-genki-bot")
	if err != nil {
		return
	}

	// read from config file
	accessToken, err := conf.GetAccessToken()
	if err != nil {
		// ignore read error
		log.Println(err)
	}
	// verify stored access token
	if accessToken != nil {
		var verified bool
		verified, err = c.verifyCredentials(accessToken)
		if err != nil {
			return
		}
		if !verified {
			log.Println("Verify credential failed.")
			accessToken = nil
		}
	}

	// get access token by "oob" process
	if accessToken == nil {
		var (
			requestToken *oauth.RequestToken
			url          string
		)
		requestToken, url, err = c.consumer.GetRequestTokenAndUrl("oob")
		if err != nil {
			return
		}
		fmt.Printf("open %s and enter PIN code.\n", url)
		fmt.Print("PIN: ")

		var verificationCode string
		if _, err = fmt.Scanln(&verificationCode); err != nil {
			return
		}
		// get credentials by PIN code
		accessToken, err = c.consumer.AuthorizeToken(requestToken, verificationCode)
		if err != nil {
			return
		}

		// verify and store to config file
		var verified bool
		verified, err = c.verifyCredentials(accessToken)
		if err != nil {
			return
		}
		if verified {
			conf.SetAccessToken(accessToken)
		} else {
			return fmt.Errorf("Verify credential failed.")
		}
	}

	c.accessToken = accessToken
	return nil
}

func (c *TwitterClient) verifyCredentials(accessToken *oauth.AccessToken) (ok bool, err error) {
	response, err := c.consumer.Get(
		"https://api.twitter.com/1.1/account/verify_credentials.json",
		map[string]string{},
		accessToken,
	)
	if err != nil {
		if err, ok := err.(oauth.HTTPExecuteError); ok {
			log.Println(err.Status)
			return false, nil
		}
		return false, err
	}
	defer response.Body.Close()

	return true, nil
}
