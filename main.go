package main

import (
	"github.com/sugyan/mentionbot"
	"log"
	"math/rand"
	"os"
	"regexp"
	"time"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	rand.Seed(time.Now().UnixNano())
}

func main() {
	config := &mentionbot.Config{
		UserID:            os.Getenv("USER_ID"),
		ConsumerKey:       os.Getenv("CONSUMER_KEY"),
		ConsumerSecret:    os.Getenv("CONSUMER_SECRET"),
		AccessToken:       os.Getenv("ACCESS_TOKEN"),
		AccessTokenSecret: os.Getenv("ACCESS_TOKEN_SECRET"),
	}
	bot := mentionbot.NewBot(config)
	bot.SetMentioner(&Genki{})
	bot.Debug(true)
	if err := bot.Run(); err != nil {
		log.Fatal(err)
	}
}

// Genki type implements mentionbot.Mentioner
type Genki struct{}

// Mention returns mention
func (*Genki) Mention(tweet *mentionbot.Tweet) (mention *string) {
	if tweet.InReplyToStatusID > 0 {
		return
	}
	if tweet.InReplyToUserID > 0 {
		return
	}
	if tweet.RetweetedStatus != nil {
		return
	}

	var shinpai string
	switch {
	case regexp.MustCompile("https?://").MatchString(tweet.Text):
		return
	case regexp.MustCompile("疲").MatchString(tweet.Text):
		if regexp.MustCompile("疲(?:れ(?:様|さ(?:ま|ん)))").MatchString(tweet.Text) {
			return
		}
		shinpai = "疲れてるの？"
	case regexp.MustCompile("凹").MatchString(tweet.Text):
		shinpai = "凹んでるの？"
	case regexp.MustCompile("心折").MatchString(tweet.Text):
		shinpai = "心折れてるの？"
	case regexp.MustCompile("(?:寂|淋)し").MatchString(tweet.Text):
		shinpai = "さびしいの？"
	case regexp.MustCompile("弱っ").MatchString(tweet.Text):
		shinpai = "弱ってるの？"
	case regexp.MustCompile("つらい").MatchString(tweet.Text):
		shinpai = "つらくても、"
	case regexp.MustCompile("死にたい").MatchString(tweet.Text):
		shinpai = "死なないで、"
	case regexp.MustCompile("お腹痛い").MatchString(tweet.Text):
		shinpai = "トイレ行って、"
	case regexp.MustCompile("(?:。。。|orz)").MatchString(tweet.Text):
	default:
		return
	}

	text := shinpai + "げんきだして！"
	return &text
}
