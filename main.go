package main

import (
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"github.com/kurrik/oauth1a"
	"github.com/kurrik/twittergo"
	"github.com/sugyan/mentionbot"
	"log"
	"math/rand"
	"net/http"
	"os"
	"regexp"
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

	bot := mentionbot.NewBot(consumerKey, consumerSecret)
	// bot.Debug(true)

	// get follwers tweets
	timeline, err := bot.FollowersTimeline(account.IdStr())
	if err != nil {
		log.Fatal(err)
	}
	for _, tweet := range timeline {
		if mention := mentionToTweet(tweet); mention != nil {
			createdAt, err := tweet.CreatedAtTime()
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("[%v] @%s (%s): %s", createdAt.Local(), tweet.User.ScreenName, tweet.User.IdStr, tweet.Text)
			log.Printf(" -> %s", mention.Text)
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

// Mention type
type Mention struct {
	StatusID int64
	Text     string
}

func mentionToTweet(tweet *anaconda.Tweet) (mention *Mention) {
	if tweet.InReplyToStatusID > 0 {
		return
	}
	if tweet.InReplyToUserID > 0 {
		return
	}
	if tweet.RetweetedStatus != nil {
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
	case regexp.MustCompile("お腹痛い").MatchString(tweet.Text):
		shinpai += "トイレ行って、"
	case regexp.MustCompile("(?:。。。|orz)").MatchString(tweet.Text):
	default:
		return
	}

	return &Mention{
		Text:     shinpai + "げんきだして！",
		StatusID: tweet.Id,
	}
}
