package slack

import (
	"testing"

	"github.com/botopolis/bot"
	"github.com/nlopes/slack"
	"github.com/stretchr/testify/assert"
)

func TestParseRoom(t *testing.T) {
	msg := slack.Message{Msg: slack.Msg{Channel: "C4321"}}
	cases := []struct {
		In  bot.Message
		Out bot.Message
	}{
		{
			In:  bot.Message{Room: "C1234"},
			Out: bot.Message{Room: "C1234"},
		},
		{
			In:  bot.Message{Room: "D1234"},
			Out: bot.Message{Room: "D1234"},
		},
		{
			In:  bot.Message{Room: "", Envelope: msg},
			Out: bot.Message{Room: "C4321", Envelope: msg},
		},
		{
			In:  bot.Message{Room: "general"},
			Out: bot.Message{Room: "C1234"},
		},
	}

	store := newTestStore()
	store.Channel.ID = "C1234"
	store.Channel.Name = "general"
	a := Adapter{Store: store}

	for _, c := range cases {
		a.parseRoom(&c.In)
		assert.Equal(t, c.Out, c.In)
	}
}

func TestParseUser(t *testing.T) {
	msg := slack.Message{Msg: slack.Msg{User: "U4321"}}
	cases := []struct {
		In  bot.Message
		Out bot.Message
	}{
		{
			In:  bot.Message{User: "U1234"},
			Out: bot.Message{User: "U1234"},
		},
		{
			In:  bot.Message{User: "", Envelope: msg},
			Out: bot.Message{User: "U4321", Envelope: msg},
		},
		{
			In:  bot.Message{User: "bob"},
			Out: bot.Message{User: "U1234"},
		},
	}

	store := newTestStore()
	store.User = slack.User{ID: "U1234", Name: "bob"}
	a := Adapter{Store: store}

	for _, c := range cases {
		a.parseUser(&c.In)
		assert.Equal(t, c.Out, c.In)
	}
}

func TestParseDM(t *testing.T) {
	cases := []struct {
		In  bot.Message
		Out bot.Message
	}{
		{
			In:  bot.Message{Room: "D1234"},
			Out: bot.Message{Room: "D1234"},
		},
		{
			In:  bot.Message{User: "U4321"},
			Out: bot.Message{User: "U4321", Room: "D1234"},
		},
	}

	store := newTestStore()
	store.IM.ID = "D1234"
	store.IM.User = "U4321"
	a := Adapter{Store: store}

	for _, c := range cases {
		a.parseDM(&c.In)
		assert.Equal(t, c.Out, c.In)
	}
}
