package main

import (
	"bytes"
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
	consumer := oauth.NewConsumer(
		"St8EoKgvUJEZ71o7GyvmSxP5w",
		"n7eS5ZDJELZLv466FYzpbMwSc9xmLkUvTkvy6AxWfvEjJZ4zWj",
		oauth.ServiceProvider{
			RequestTokenUrl:   "https://api.twitter.com/oauth/request_token",
			AuthorizeTokenUrl: "https://api.twitter.com/oauth/authorize",
			AccessTokenUrl:    "https://api.twitter.com/oauth/access_token",
		},
	)
	requestToken, url, err := consumer.GetRequestTokenAndUrl("oob")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("open %s\n", url)
	fmt.Print("PIN: ")

	var verificationCode string
	if _, err = fmt.Scanln(&verificationCode); err != nil {
		log.Fatal(err)
	}

	accessToken, err := consumer.AuthorizeToken(requestToken, verificationCode)
	if err != nil {
		log.Fatal(err)
	}

	return &TwitterClient{
		consumer:    consumer,
		accessToken: accessToken,
	}
}

// VerifyCredentials returns a representation of the requesting user
func (c *TwitterClient) VerifyCredentials() string {
	response, err := c.consumer.Get(
		"https://api.twitter.com/1.1/account/verify_credentials.json",
		map[string]string{},
		c.accessToken,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	buffer := new(bytes.Buffer)
	if _, err := buffer.ReadFrom(response.Body); err != nil {
		log.Fatal(err)
	}
	return buffer.String()
}
