package main

import (
	"github.com/kurrik/oauth1a"
	"github.com/kurrik/twittergo"
	"log"
	"net/http"
	"os"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	clientConfig := &oauth1a.ClientConfig{
		ConsumerKey:    os.Getenv("CONSUMER_KEY"),
		ConsumerSecret: os.Getenv("CONSUMER_SECRET"),
	}
	userConfig := oauth1a.NewAuthorizedConfig(os.Getenv("ACCESS_TOKEN"), os.Getenv("ACCESS_TOKEN_SECRET"))
	client := twittergo.NewClient(clientConfig, userConfig)

	req, err := http.NewRequest("GET", "/1.1/account/verify_credentials.json", nil)
	if err != nil {
		log.Fatal(err)
	}
	resp, err := client.SendRequest(req)
	if err != nil {
		log.Fatal(err)
	}

	user := &twittergo.User{}
	resp.Parse(user)
	log.Println(user)
}
