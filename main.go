package main

import (
	"github.com/kurrik/oauth1a"
	"github.com/kurrik/twittergo"
	"github.com/sugyan/replybot"
	"log"
	"net/http"
	"os"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	var (
		consumerKey       = os.Getenv("CONSUMER_KEY")
		consumerSecret    = os.Getenv("CONSUMER_SECRET")
		accessToken       = os.Getenv("ACCESS_TOKEN")
		accessTokenSecret = os.Getenv("ACCESS_TOKEN_SECRET")
	)
	account, err := getAccount(consumerKey, consumerSecret, accessToken, accessTokenSecret)
	if err != nil {
		log.Fatal(err)
	}

	bot := replybot.NewBot(consumerKey, consumerSecret)
	ids, err := bot.FollowersIDs(account.IdStr())
	if err != nil {
		log.Fatal(err)
	}
	log.Println(len(ids))
}

func getAccount(consumerKey string, consumerSecret string, accessToken string, accessTokenSecret string) (*twittergo.User, error) {
	clientConfig := &oauth1a.ClientConfig{
		ConsumerKey:    consumerKey,
		ConsumerSecret: consumerSecret,
	}
	userConfig := oauth1a.NewAuthorizedConfig(accessToken, accessTokenSecret)
	client := twittergo.NewClient(clientConfig, userConfig)

	req, err := http.NewRequest("GET", "/1.1/account/verify_credentials.json", nil)
	if err != nil {
		return nil, err
	}
	res, err := client.SendRequest(req)
	if err != nil {
		return nil, err
	}
	user := &twittergo.User{}
	res.Parse(user)

	return user, nil
}
