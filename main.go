package main

import (
	"github.com/kurrik/oauth1a"
	"github.com/kurrik/twittergo"
	"log"
	"net/http"
	"net/url"
	"os"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	bot, err := getBotAccount()
	if err != nil {
		log.Fatal(err)
	}
	clientConfig := &oauth1a.ClientConfig{
		ConsumerKey:    os.Getenv("CONSUMER_KEY"),
		ConsumerSecret: os.Getenv("CONSUMER_SECRET"),
	}
	client := twittergo.NewClient(clientConfig, nil)

	var count int
	var cursor string
	for {
		query := url.Values{}
		query.Set("user_id", bot.IdStr())
		query.Set("stringify_ids", "true")
		query.Set("count", "5000")
		if cursor != "" {
			query.Set("cursor", cursor)
		}
		req, err := http.NewRequest("GET", "/1.1/followers/ids.json?"+query.Encode(), nil)
		if err != nil {
			log.Fatal(err)
		}
		res, err := client.SendRequest(req)
		if err != nil {
			log.Fatal(err)
		}

		results := &CursoredIDs{}
		if err := res.Parse(results); err != nil {
			log.Fatal(err)
		}

		log.Println(results)
		for _, id := range results.IDs() {
			count++
			log.Println(id)
		}

		log.Printf("%v, %v", results.PreviousCursorStr(), results.NextCursorStr())
		if res.HasRateLimit() {
			log.Printf("Rate Limit: %v / %v, (next reset: %v)\n", res.RateLimitRemaining(), res.RateLimit(), res.RateLimitReset())
		}

		if results.NextCursorStr() == "0" {
			break
		} else {
			cursor = results.NextCursorStr()
		}
	}
	log.Printf("%v users", count)
}

func getBotAccount() (user *twittergo.User, err error) {
	clientConfig := &oauth1a.ClientConfig{
		ConsumerKey:    os.Getenv("CONSUMER_KEY"),
		ConsumerSecret: os.Getenv("CONSUMER_SECRET"),
	}
	userConfig := oauth1a.NewAuthorizedConfig(os.Getenv("ACCESS_TOKEN"), os.Getenv("ACCESS_TOKEN_SECRET"))
	client := twittergo.NewClient(clientConfig, userConfig)

	req, err := http.NewRequest("GET", "/1.1/account/verify_credentials.json", nil)
	if err != nil {
		return nil, err
	}
	res, err := client.SendRequest(req)
	if err != nil {
		return nil, err
	}
	user = &twittergo.User{}
	res.Parse(user)

	return user, nil
}
