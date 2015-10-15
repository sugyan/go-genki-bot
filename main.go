package main

import (
	"github.com/kurrik/oauth1a"
	"github.com/kurrik/twittergo"
	"github.com/sugyan/mentionbot"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	rand.Seed(time.Now().UnixNano())
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

	// get follwers ids
	bot := mentionbot.NewBot(consumerKey, consumerSecret)
	// bot.Debug(true)
	ids, err := bot.FollowersIDs(account.IdStr())
	if err != nil {
		log.Fatal(err)
	}
	// shuffle
	for i := len(ids) - 1; i >= 0; i-- {
		j := rand.Intn(i + 1)
		ids[i], ids[j] = ids[j], ids[i]
	}
	// get followers list
	users, err := bot.UsersLookup(ids[0:100])
	if err != nil {
		log.Fatal(err)
	}
	for _, user := range users {
		status := user.Status()
		if status != nil {
			log.Printf("[%v] @%s: %s", status.CreatedAt().Local(), user.ScreenName(), status.Text())
		}
	}
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
