package main

import (
	"fmt"
	"log"
	"regexp"
)

// GenkiBot type
type GenkiBot struct {
	client *TwitterClient
}

// Mention type
type Mention struct {
	StatusID int64
	Text     string
}

// NewGenkiBot returns new GenkiBot instance
func NewGenkiBot() *GenkiBot {
	return &GenkiBot{
		client: NewTwitterClient(),
	}
}

// MentionToTweet returns mention to tweet
func (*GenkiBot) MentionToTweet(tweet *Tweet) (mention *Mention) {
	if tweet.InReplyToStatusID > 0 {
		return
	}
	if tweet.InReplyToUserID > 0 {
		return
	}
	if tweet.RetweetedStatus.ID > 0 {
		return
	}

	shinpai := fmt.Sprintf("@%s ", tweet.User.ScreenName)
	switch {
	case regexp.MustCompile("https?://").MatchString(tweet.Text):
		return
	case regexp.MustCompile("疲").MatchString(tweet.Text):
		if regexp.MustCompile("疲(?:れ(?:様|さ(?:ま|ん)))").MatchString(tweet.Text) {
			return
		}
		shinpai += "疲れてるの？"
	case regexp.MustCompile("凹").MatchString(tweet.Text):
		shinpai += "凹んでるの？"
	case regexp.MustCompile("心折").MatchString(tweet.Text):
		shinpai += "心折れてるの？"
	case regexp.MustCompile("(?:寂|淋)し").MatchString(tweet.Text):
		shinpai += "さびしいの？"
	case regexp.MustCompile("弱っ").MatchString(tweet.Text):
		shinpai += "弱ってるの？"
	case regexp.MustCompile("つらい").MatchString(tweet.Text):
		shinpai += "つらくても、"
	case regexp.MustCompile("死にたい").MatchString(tweet.Text):
		shinpai += "死なないで、"
	case regexp.MustCompile("(?:。。。|orz)").MatchString(tweet.Text):
	default:
		return
	}

	return &Mention{
		Text:     shinpai + "げんきだして！",
		StatusID: tweet.ID,
	}
}

// Run starts genki-bot
func (bot *GenkiBot) Run(user *string) {
	if err := bot.client.PrepareAccessToken(user); err != nil {
		log.Fatal("error: ", err)
	}
	stream, err := bot.client.UserStream()
	if err != nil {
		log.Fatal("error: ", err)
	}

	log.Println("userstream started.")
	for {
		tweet, err := stream.NextTweet()
		if err != nil {
			log.Fatal(err)
		}
		log.Println(tweet)

		mention := bot.MentionToTweet(tweet)
		if mention != nil {
			log.Println(mention)
		}
	}
}
