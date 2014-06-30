package main

import (
	"log"
)

func main() {
	client := NewTwitterClient()
	if err := client.PrepareAccessToken(); err != nil {
		log.Fatal("error: ", err)
	}

	stream, err := client.UserStream()
	if err != nil {
		log.Fatal("error: ", err)
	}

	for {
		tweet, err := stream.NextTweet()
		if err != nil {
			log.Fatal(err)
		}
		log.Println(tweet)
	}
}
