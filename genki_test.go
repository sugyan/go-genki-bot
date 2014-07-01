package main

import (
	"testing"
)

func TestGenki(t *testing.T) {
	var result *Mention

	// empty tweet
	result = (&Tweet{}).Genki()
	if result != nil {
		t.Error("should not return mention")
	}

	// normal tweet
	result = (&Tweet{
		Text: "とても普通の健全なツイート",
	}).Genki()
	if result != nil {
		t.Error("should not return mention")
	}

	// 疲れた
	result = (&Tweet{
		ID:   555,
		Text: "はー疲れた",
		User: TweetUser{
			ScreenName: "hoge",
		},
	}).Genki()
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
		result = (&Tweet{
			Text: text,
		}).Genki()
		if result != nil {
			t.Error("should not return mention")
		}
	}

	// 凹
	result = (&Tweet{
		Text: "凹んでる",
		User: TweetUser{
			ScreenName: "hoge",
		},
	}).Genki()
	if result == nil {
		t.Fatal("should return mention")
	}
	if result.Text != "@hoge 凹んでるの？げんきだして！" {
		t.Errorf("incorrect text: %s", result.Text)
	}

	// 心折
	result = (&Tweet{
		Text: "心折れた。",
		User: TweetUser{
			ScreenName: "hoge",
		},
	}).Genki()
	if result == nil {
		t.Fatal("should return mention")
	}
	if result.Text != "@hoge 心折れてるの？げんきだして！" {
		t.Errorf("incorrect text: %s", result.Text)
	}

	// さびしい
	for _, text := range []string{"寂しい", "淋しい"} {
		result = (&Tweet{
			Text: text,
			User: TweetUser{
				ScreenName: "hoge",
			},
		}).Genki()
		if result == nil {
			t.Fatal("should return mention")
		}
		if result.Text != "@hoge さびしいの？げんきだして！" {
			t.Errorf("incorrect text: %s", result.Text)
		}
	}

	// 弱ってる
	result = (&Tweet{
		Text: "弱ったなぁ",
		User: TweetUser{
			ScreenName: "hoge",
		},
	}).Genki()
	if result == nil {
		t.Fatal("should return mention")
	}
	if result.Text != "@hoge 弱ってるの？げんきだして！" {
		t.Errorf("incorrect text: %s", result.Text)
	}

	// つらい
	result = (&Tweet{
		Text: "とてもつらい",
		User: TweetUser{
			ScreenName: "hoge",
		},
	}).Genki()
	if result == nil {
		t.Fatal("should return mention")
	}
	if result.Text != "@hoge つらくても、げんきだして！" {
		t.Errorf("incorrect text: %s", result.Text)
	}

	// 死にたい
	result = (&Tweet{
		Text: "死にたい〜",
		User: TweetUser{
			ScreenName: "hoge",
		},
	}).Genki()
	if result == nil {
		t.Fatal("should return mention")
	}
	if result.Text != "@hoge 死なないで、げんきだして！" {
		t.Errorf("incorrect text: %s", result.Text)
	}

	// その他
	for _, text := range []string{"うわ。。。", "ミスった orz"} {
		result = (&Tweet{
			Text: text,
			User: TweetUser{
				ScreenName: "hoge",
			},
		}).Genki()
		if result == nil {
			t.Fatal("should return mention")
		}
		if result.Text != "@hoge げんきだして！" {
			t.Errorf("incorrect text: %s", result.Text)
		}
	}

	// include URL
	result = (&Tweet{
		Text: "疲れてるしつらい。 http://example.com",
	}).Genki()
	if result != nil {
		t.Error("should not return mention")
	}

	// reply, retweet
	result = (&Tweet{
		Text:              "疲れてるしつらい。",
		InReplyToStatusID: 42,
	}).Genki()
	if result != nil {
		t.Error("should not return mention")
	}

	result = (&Tweet{
		Text: "疲れてるしつらい。",
		RetweetedStatus: RetweetedStatus{
			ID: 42,
		},
	}).Genki()
	if result != nil {
		t.Error("should not return mention")
	}
}
