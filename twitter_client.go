package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/mrjones/oauth"
	"log"
	"strconv"
)

// TwitterClient type
type TwitterClient struct {
	consumer    *oauth.Consumer
	accessToken *oauth.AccessToken
}

// TweetUser type
type TweetUser struct {
	ID         uint64
	Name       string
	ScreenName string `json:"screen_name"`
}

// RetweetedStatus type
type RetweetedStatus struct {
	ID        uint64
	Text      string
	CreatedAt string `json:"created_at"`
	User      TweetUser
}

// Tweet type
type Tweet struct {
	ID                uint64
	Text              string
	CreatedAt         string          `json:"created_at"`
	InReplyToUserID   uint64          `json:"in_reply_to_user_id"`
	InReplyToStatusID uint64          `json:"in_reply_to_status_id"`
	RetweetedStatus   RetweetedStatus `json:"retweeted_status"`
	User              TweetUser
}

// Stream type
type Stream struct {
	scanner *bufio.Scanner
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
func (c *TwitterClient) PrepareAccessToken(user *string) (err error) {
	conf, err := NewConfig("go-genki-bot")
	if err != nil {
		return
	}

	// read from config file
	accessToken, err := conf.GetAccessToken(user)
	if err != nil {
		// ignore read error
		log.Println(err)
	}
	// verify stored access token
	if accessToken != nil {
		var verified bool
		verified, err = c.verifyCredentials(accessToken.Token, accessToken.Secret)
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
		requestToken, url, err := c.consumer.GetRequestTokenAndUrl("oob")
		if err != nil {
			return err
		}
		fmt.Printf("open %s and enter PIN code.\n", url)
		fmt.Print("PIN: ")

		var verificationCode string
		if _, err := fmt.Scanln(&verificationCode); err != nil {
			return err
		}
		// get credentials by PIN code
		verifiedToken, err := c.consumer.AuthorizeToken(requestToken, verificationCode)
		if err != nil {
			return err
		}

		// verify and store to config file
		verified, err := c.verifyCredentials(verifiedToken.Token, verifiedToken.Secret)
		if err != nil {
			return err
		}

		accessToken = &AccessToken{
			Token:  verifiedToken.Token,
			Secret: verifiedToken.Secret,
		}
		if verified {
			if name, ok := verifiedToken.AdditionalData["screen_name"]; ok {
				conf.SetAccessToken(name, accessToken)
			}
			conf.SetAccessToken("default", accessToken)
		} else {
			return fmt.Errorf("Verify credential failed.")
		}
	}

	c.accessToken = &oauth.AccessToken{
		Token:  accessToken.Token,
		Secret: accessToken.Secret,
	}
	return nil
}

// UserStream returns user stream
func (c *TwitterClient) UserStream() (stream *Stream, err error) {
	resp, err := c.consumer.Get(
		"https://userstream.twitter.com/1.1/user.json",
		map[string]string{},
		c.accessToken,
	)
	if err != nil {
		return
	}

	return &Stream{
		scanner: bufio.NewScanner(resp.Body),
	}, nil
}

// Mention updates status with in_reply_to_status_id
func (c *TwitterClient) Mention(mention *Mention) (tweet *Tweet, err error) {
	resp, err := c.consumer.Post(
		"https://api.twitter.com/1.1/statuses/update.json",
		map[string]string{
			"status":                mention.Text,
			"in_reply_to_status_id": strconv.FormatUint(mention.StatusID, 10),
		},
		c.accessToken,
	)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	json.NewDecoder(resp.Body).Decode(&tweet)
	return tweet, nil
}

func (c *TwitterClient) verifyCredentials(token string, secret string) (ok bool, err error) {
	resp, err := c.consumer.Get(
		"https://api.twitter.com/1.1/account/verify_credentials.json",
		map[string]string{},
		&oauth.AccessToken{
			Token:  token,
			Secret: secret,
		},
	)
	if err != nil {
		if err, ok := err.(oauth.HTTPExecuteError); ok {
			log.Println(err.Status)
			return false, nil
		}
		return false, err
	}
	defer resp.Body.Close()

	return true, nil
}

// NextTweet returns tweet from stream
func (s *Stream) NextTweet() (tweet *Tweet, err error) {
	for s.scanner.Err() == nil {
		var bytes []byte
		bytes, err = func() ([]byte, error) {
			for {
				if !s.scanner.Scan() {
					return nil, s.scanner.Err()
				}
				bytes := s.scanner.Bytes()
				if len(bytes) > 0 {
					return bytes, nil
				}
			}
		}()
		if err != nil {
			return nil, err
		}
		if err = json.Unmarshal(bytes, &tweet); err != nil {
			return nil, err
		}
		// TODO?
		if tweet.ID > 0 {
			return tweet, nil
		}
	}
	return nil, s.scanner.Err()
}
