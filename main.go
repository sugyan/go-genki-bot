package main

import (
	"github.com/ChimeraCoder/anaconda"
	"log"
	"net/url"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	// TODO: config...
	anaconda.SetConsumerKey("")
	anaconda.SetConsumerSecret("")
	api := anaconda.NewTwitterApi("", "")
	v := url.Values{}
	user, err := api.GetSelf(v)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(user)
}
