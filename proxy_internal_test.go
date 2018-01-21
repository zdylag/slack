package slack

import (
	"testing"

	"github.com/berfarah/gobot"
	"github.com/nlopes/slack"
	"github.com/stretchr/testify/assert"
)

func TestProxy(t *testing.T) {
	assert := assert.New(t)

	room := "C1234"
	user := "U1234"
	text := "foo bar"
	topic := "generally awesome"

	proxyTestCases := []struct {
		In     slack.MessageEvent
		Assert func(gobot.Message, gobot.Message)
	}{
		{
			In: slack.MessageEvent(slack.Message{Msg: slack.Msg{
				SubType: "",
				User:    user,
				Channel: room,
				Text:    text,
			}}),
			Assert: func(expected gobot.Message, result gobot.Message) {
				assert.Equal(gobot.MessageEvent, result.Event)
				assert.Equal(expected.Text, result.Text)
			},
		},
		{
			In: slack.MessageEvent(slack.Message{Msg: slack.Msg{
				SubType: "channel_join",
				User:    user,
				Channel: room,
				Text:    text,
			}}),
			Assert: func(expected gobot.Message, result gobot.Message) {
				assert.Equal(gobot.EnterEvent, result.Event)
				assert.Equal(expected.Text, result.Text)
			},
		},
		{
			In: slack.MessageEvent(slack.Message{Msg: slack.Msg{
				SubType: "channel_leave",
				User:    user,
				Channel: room,
				Text:    text,
			}}),
			Assert: func(expected gobot.Message, result gobot.Message) {
				assert.Equal(gobot.LeaveEvent, result.Event)
				assert.Equal(expected.Text, result.Text)
			},
		},
		{
			In: slack.MessageEvent(slack.Message{Msg: slack.Msg{
				SubType: "channel_topic",
				User:    user,
				Channel: room,
				Topic:   topic,
				Text:    text,
			}}),
			Assert: func(expected gobot.Message, result gobot.Message) {
				assert.Equal(gobot.TopicEvent, result.Event)
				assert.Equal(expected.Topic, result.Topic)
			},
		},
	}

	p := proxy{
		Adapter:   New(""),
		formatter: formatter{testStore{}},
	}
	for _, c := range proxyTestCases {
		in := make(chan slack.RTMEvent, 2)
		out := make(chan gobot.Message, 2)
		go p.Forward(in, out)

		in <- slack.RTMEvent{Type: "message", Data: &c.In}
		c.Assert(gobot.Message{
			User:     "bob",
			Room:     "general",
			Text:     text,
			Envelope: &c.In,
			Topic:    topic,
		}, <-out)
		close(in)
	}
}
