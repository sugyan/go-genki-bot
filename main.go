package main

import (
	"log"
)

func main() {
	client := NewTwitterClient()
	log.Println(client.VerifyCredentials())
}
