package main

import (
	"log"
)

func main() {

	client := NewTwitterClient()
	if err := client.PrepareAccessToken(); err != nil {
		log.Fatal("error: ", err)
	}
	log.Println(client)
}
