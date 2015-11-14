package mentionbot

import (
	"github.com/garyburd/go-oauth/oauth"
	"log"
	"sort"
	"sync"
	"time"
)

// Mentioner interface
type Mentioner interface {
	Mention(*Tweet) *string
}

// Bot type
type Bot struct {
	userID      string
	client      *oauth.Client
	credentials *oauth.Credentials
	mentioner   Mentioner
	idsStore    *idsStore
	apiBase     string
	debug       bool
}

// Config type
type Config struct {
	UserID            string
	ConsumerKey       string
	ConsumerSecret    string
	AccessToken       string
	AccessTokenSecret string
}

// NewBot returns new bot
func NewBot(config *Config) *Bot {
	return &Bot{
		userID: config.UserID,
		client: &oauth.Client{
			Credentials: oauth.Credentials{
				Token:  config.ConsumerKey,
				Secret: config.ConsumerSecret,
			},
		},
		credentials: &oauth.Credentials{
			Token:  config.AccessToken,
			Secret: config.AccessTokenSecret,
		},
		idsStore: &idsStore{},
		apiBase:  "https://api.twitter.com/1.1",
	}
}

// Debug sets debug flag
func (bot *Bot) Debug(enabled bool) {
	bot.debug = enabled
}

// SetMentioner sets mentioner instance
func (bot *Bot) SetMentioner(m Mentioner) {
	bot.mentioner = m
}

// Run bot
func (bot *Bot) Run() (err error) {
	rateLimitStatusResult, err := bot.rateLimitStatus([]string{"users"})
	if err != nil {
		return err
	}
	latestRateLimit := rateLimitStatusResult.results.(rateLimitStatusResources).Users["/users/lookup"]
	latestCreatedAt := time.Now().Add(-15 * time.Minute)

	for {
		// get follwers tweets
		timeline, rateLimit, err := bot.followersTimeline(bot.userID, latestCreatedAt)
		if err != nil {
			return err
		}

		if bot.debug {
			log.Printf("%d tweets fetched", len(timeline))
		}
		for _, tweet := range timeline {
			createdAt, err := tweet.CreatedAtTime()
			if err != nil {
				return err
			}
			if bot.mentioner != nil {
				mention := bot.mentioner.Mention(tweet)
				if mention == nil {
					continue
				}
				if bot.debug {
					log.Printf("(%s)[%v] @%s: %s", tweet.IDStr, createdAt.Local(), tweet.User.ScreenName, tweet.Text)
				}
				updated, err := bot.statusesUpdate(*mention, tweet)
				if err != nil {
					return nil
				}
				log.Println(updated.results.(Tweet).Text)
			}
		}
		// udpate latestCreatedAt
		if len(timeline) > 0 {
			latestCreatedAt, err = timeline[len(timeline)-1].CreatedAtTime()
			if err != nil {
				return err
			}
		}

		// calculate waiting time
		wait := rateLimit.waitSeconds(&latestRateLimit)
		// update latestRateLimit
		latestRateLimit = *rateLimit

		if bot.debug {
			log.Printf("wait %d seconds for next loop", wait)
		}
		<-time.Tick(time.Second * time.Duration(wait))
	}
}

func (bot *Bot) followersTimeline(userID string, since time.Time) (timeline timeline, rateLimit *rateLimitStatus, err error) {
	defer func() {
		// sort by createdAt
		if timeline != nil {
			sort.Sort(timeline)
		}
	}()

	// IDs from cache or API
	ids := bot.idsStore.pickIds()
	if ids == nil {
		idsResults, err := bot.followersIDs(userID)
		if err != nil {
			return nil, nil, err
		}
		results := idsResults.results.([]int64)
		bot.idsStore.setIds(results, 15*time.Minute)
		ids = bot.idsStore.pickIds()
	}

	type result struct {
		apiResult *apiResult
		err       error
	}
	cancel := make(chan struct{})
	defer close(cancel)

	in := make(chan []int64)
	out := make(chan result)
	// input ids (user ids length upto 100)
	// TODO: shuffle ids?
	go func() {
		for m := 0; ; m += 100 {
			n := m + 100
			if n > len(ids) {
				n = len(ids)
			}
			if n-m < 1 {
				break
			}
			in <- ids[m:n]
		}
		close(in)
	}()
	// parallelize request (bounding the number of workers)
	const numWorkers = 5
	wg := sync.WaitGroup{}
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for ids := range in {
				results, err := bot.usersLookup(ids)
				select {
				case out <- result{apiResult: results, err: err}:
				case <-cancel:
					return
				}
			}
		}()
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	// collect all results
	rateLimit = &rateLimitStatus{}
Loop:
	for {
		select {
		case result, ok := <-out:
			if !ok {
				break Loop
			}
			if result.err != nil {
				return nil, nil, result.err
			}
			apiResult := result.apiResult
			if apiResult.rateLimit != nil {
				if (apiResult.rateLimit.Reset > rateLimit.Reset) || (apiResult.rateLimit.Remaining < rateLimit.Remaining) {
					rateLimit = apiResult.rateLimit
				}
			}
			// make results
			for _, user := range apiResult.results.([]User) {
				tweet := user.Status
				if tweet != nil {
					createdAtTime, err := tweet.CreatedAtTime()
					if err != nil {
						return nil, nil, err
					}
					if createdAtTime.After(since) {
						tweet.User = user
						timeline = append(timeline, tweet)
					}
				}
			}
		}
	}
	return
}
