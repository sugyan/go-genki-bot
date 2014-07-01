package main

import (
	"fmt"
	"regexp"
)

// Mention type
type Mention struct {
	StatusID int64
	Text     string
}

// Genki returns mention data to send
func (t *Tweet) Genki() (m *Mention) {
	if t.InReplyToStatusID > 0 {
		return
	}
	if t.InReplyToUserID > 0 {
		return
	}
	if t.RetweetedStatus.ID > 0 {
		return
	}

	shinpai := fmt.Sprintf("@%s ", t.User.ScreenName)
	switch {
	case regexp.MustCompile("https?://").MatchString(t.Text):
		return
	case regexp.MustCompile("疲").MatchString(t.Text):
		if regexp.MustCompile("疲(?:れ(?:様|さ(?:ま|ん)))").MatchString(t.Text) {
			return
		}
		shinpai += "疲れてるの？"
	case regexp.MustCompile("凹").MatchString(t.Text):
		shinpai += "凹んでるの？"
	case regexp.MustCompile("心折").MatchString(t.Text):
		shinpai += "心折れてるの？"
	case regexp.MustCompile("(?:寂|淋)し").MatchString(t.Text):
		shinpai += "さびしいの？"
	case regexp.MustCompile("弱っ").MatchString(t.Text):
		shinpai += "弱ってるの？"
	case regexp.MustCompile("つらい").MatchString(t.Text):
		shinpai += "つらくても、"
	case regexp.MustCompile("死にたい").MatchString(t.Text):
		shinpai += "死なないで、"
	case regexp.MustCompile("(?:。。。|orz)").MatchString(t.Text):
	default:
		return
	}

	return &Mention{
		Text:     shinpai + "げんきだして！",
		StatusID: t.ID,
	}
}
