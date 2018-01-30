package slack

import (
	"testing"

	"github.com/berfarah/gobot"
	"github.com/nlopes/slack"
	"github.com/stretchr/testify/assert"
)

func TestParseRoom(t *testing.T) {
	msg := slack.Message{Msg: slack.Msg{Channel: "C4321"}}
	cases := []struct {
		In  gobot.Message
		Out gobot.Message
	}{
		{
			In:  gobot.Message{Room: "C1234"},
			Out: gobot.Message{Room: "C1234"},
		},
		{
			In:  gobot.Message{Room: "D1234"},
			Out: gobot.Message{Room: "D1234"},
		},
		{
			In:  gobot.Message{Room: "", Envelope: msg},
			Out: gobot.Message{Room: "C4321", Envelope: msg},
		},
		{
			In:  gobot.Message{Room: "general"},
			Out: gobot.Message{Room: "C1234"},
		},
	}

	a := Adapter{Store: testStore{}}

	for _, c := range cases {
		a.parseRoom(&c.In)
		assert.Equal(t, c.Out, c.In)
	}
}

func TestParseUser(t *testing.T) {
	msg := slack.Message{Msg: slack.Msg{User: "U4321"}}
	cases := []struct {
		In  gobot.Message
		Out gobot.Message
	}{
		{
			In:  gobot.Message{User: "U1234"},
			Out: gobot.Message{User: "U1234"},
		},
		{
			In:  gobot.Message{User: "", Envelope: msg},
			Out: gobot.Message{User: "U4321", Envelope: msg},
		},
		{
			In:  gobot.Message{User: "bob"},
			Out: gobot.Message{User: "U1234"},
		},
	}

	a := Adapter{Store: testStore{}}

	for _, c := range cases {
		a.parseUser(&c.In)
		assert.Equal(t, c.Out, c.In)
	}
}

func TestParseDM(t *testing.T) {
	cases := []struct {
		In  gobot.Message
		Out gobot.Message
	}{
		{
			In:  gobot.Message{Room: "D1234"},
			Out: gobot.Message{Room: "D1234"},
		},
		{
			In:  gobot.Message{User: "U4321"},
			Out: gobot.Message{User: "U4321", Room: "D1234"},
		},
	}

	a := Adapter{Store: testStore{}}

	for _, c := range cases {
		a.parseDM(&c.In)
		assert.Equal(t, c.Out, c.In)
	}
}
