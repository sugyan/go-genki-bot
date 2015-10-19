package main

import (
	"github.com/sugyan/mentionbot"
	"log"
	"math/rand"
	"os"
	"time"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	rand.Seed(time.Now().UnixNano())
}

func main() {
	var (
		botUserID      = os.Getenv("USER_ID")
		consumerKey    = os.Getenv("CONSUMER_KEY")
		consumerSecret = os.Getenv("CONSUMER_SECRET")
	)
	bot := mentionbot.NewBot(botUserID, consumerKey, consumerSecret)
	bot.Debug(true)
	if err := bot.Run(); err != nil {
		log.Fatal(err)
	}
}

// // Mention type
// type Mention struct {
// 	StatusID int64
// 	Text     string
// }

// func mentionToTweet(tweet *anaconda.Tweet) (mention *Mention) {
// 	if tweet.InReplyToStatusID > 0 {
// 		return
// 	}
// 	if tweet.InReplyToUserID > 0 {
// 		return
// 	}
// 	if tweet.RetweetedStatus != nil {
// 		return
// 	}

// 	shinpai := fmt.Sprintf("@%s ", tweet.User.ScreenName)
// 	switch {
// 	case regexp.MustCompile("https?://").MatchString(tweet.Text):
// 		return
// 	case regexp.MustCompile("疲").MatchString(tweet.Text):
// 		if regexp.MustCompile("疲(?:れ(?:様|さ(?:ま|ん)))").MatchString(tweet.Text) {
// 			return
// 		}
// 		shinpai += "疲れてるの？"
// 	case regexp.MustCompile("凹").MatchString(tweet.Text):
// 		shinpai += "凹んでるの？"
// 	case regexp.MustCompile("心折").MatchString(tweet.Text):
// 		shinpai += "心折れてるの？"
// 	case regexp.MustCompile("(?:寂|淋)し").MatchString(tweet.Text):
// 		shinpai += "さびしいの？"
// 	case regexp.MustCompile("弱っ").MatchString(tweet.Text):
// 		shinpai += "弱ってるの？"
// 	case regexp.MustCompile("つらい").MatchString(tweet.Text):
// 		shinpai += "つらくても、"
// 	case regexp.MustCompile("死にたい").MatchString(tweet.Text):
// 		shinpai += "死なないで、"
// 	case regexp.MustCompile("お腹痛い").MatchString(tweet.Text):
// 		shinpai += "トイレ行って、"
// 	case regexp.MustCompile("(?:。。。|orz)").MatchString(tweet.Text):
// 	default:
// 		return
// 	}

// 	return &Mention{
// 		Text:     shinpai + "げんきだして！",
// 		StatusID: tweet.Id,
// 	}
// }
