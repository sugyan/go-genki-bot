package main

import (
	"testing"
)

func TestMention(t *testing.T) {
	var result *Mention
	bot := NewGenkiBot()

	// empty tweet
	result = bot.MentionToTweet(&Tweet{})
	if result != nil {
		t.Error("should not return mention")
	}

	// normal tweet
	result = bot.MentionToTweet(&Tweet{
		Text: "とても普通の健全なツイート",
	})
	if result != nil {
		t.Error("should not return mention")
	}

	// 疲れた
	result = bot.MentionToTweet(&Tweet{
		ID:   555,
		Text: "はー疲れた",
		User: TweetUser{
			ScreenName: "hoge",
		},
	})
	if result == nil {
		t.Fatal("should return mention")
	}
	if result.Text != "@hoge 疲れてるの？げんきだして！" {
		t.Errorf("incorrect text: %s", result.Text)
	}
	if result.StatusID != 555 {
		t.Errorf("incorrect status id: %s", result.StatusID)
	}
	// 疲れてない
	for _, text := range []string{"お疲れさまでした。", "お疲れさん", "お疲れ様です"} {
		result = bot.MentionToTweet(&Tweet{
			Text: text,
		})
		if result != nil {
			t.Error("should not return mention")
		}
	}

	// 凹
	result = bot.MentionToTweet(&Tweet{
		Text: "凹んでる",
		User: TweetUser{
			ScreenName: "hoge",
		},
	})
	if result == nil {
		t.Fatal("should return mention")
	}
	if result.Text != "@hoge 凹んでるの？げんきだして！" {
		t.Errorf("incorrect text: %s", result.Text)
	}

	// 心折
	result = bot.MentionToTweet(&Tweet{
		Text: "心折れた。",
		User: TweetUser{
			ScreenName: "hoge",
		},
	})
	if result == nil {
		t.Fatal("should return mention")
	}
	if result.Text != "@hoge 心折れてるの？げんきだして！" {
		t.Errorf("incorrect text: %s", result.Text)
	}

	// さびしい
	for _, text := range []string{"寂しい", "淋しい"} {
		result = bot.MentionToTweet(&Tweet{
			Text: text,
			User: TweetUser{
				ScreenName: "hoge",
			},
		})
		if result == nil {
			t.Fatal("should return mention")
		}
		if result.Text != "@hoge さびしいの？げんきだして！" {
			t.Errorf("incorrect text: %s", result.Text)
		}
	}

	// 弱ってる
	result = bot.MentionToTweet(&Tweet{
		Text: "弱ったなぁ",
		User: TweetUser{
			ScreenName: "hoge",
		},
	})
	if result == nil {
		t.Fatal("should return mention")
	}
	if result.Text != "@hoge 弱ってるの？げんきだして！" {
		t.Errorf("incorrect text: %s", result.Text)
	}

	// つらい
	result = bot.MentionToTweet(&Tweet{
		Text: "とてもつらい",
		User: TweetUser{
			ScreenName: "hoge",
		},
	})
	if result == nil {
		t.Fatal("should return mention")
	}
	if result.Text != "@hoge つらくても、げんきだして！" {
		t.Errorf("incorrect text: %s", result.Text)
	}

	// 死にたい
	result = bot.MentionToTweet(&Tweet{
		Text: "死にたい〜",
		User: TweetUser{
			ScreenName: "hoge",
		},
	})
	if result == nil {
		t.Fatal("should return mention")
	}
	if result.Text != "@hoge 死なないで、げんきだして！" {
		t.Errorf("incorrect text: %s", result.Text)
	}

	// その他
	for _, text := range []string{"うわ。。。", "ミスった orz"} {
		result = bot.MentionToTweet(&Tweet{
			Text: text,
			User: TweetUser{
				ScreenName: "hoge",
			},
		})
		if result == nil {
			t.Fatal("should return mention")
		}
		if result.Text != "@hoge げんきだして！" {
			t.Errorf("incorrect text: %s", result.Text)
		}
	}

	// include URL
	result = bot.MentionToTweet(&Tweet{
		Text: "疲れてるしつらい。 http://example.com",
	})
	if result != nil {
		t.Error("should not return mention")
	}

	// reply, retweet
	result = bot.MentionToTweet(&Tweet{
		Text:              "疲れてるしつらい。",
		InReplyToStatusID: 42,
	})
	if result != nil {
		t.Error("should not return mention")
	}

	result = bot.MentionToTweet(&Tweet{
		Text: "疲れてるしつらい。",
		RetweetedStatus: RetweetedStatus{
			ID: 42,
		},
	})
	if result != nil {
		t.Error("should not return mention")
	}
}
