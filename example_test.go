package slack_test

import (
	"fmt"
	"os"

	"github.com/berfarah/gobot"
	"github.com/berfarah/gobot-slack"
	slacker "github.com/nlopes/slack"
)

func Example_basic() {
	robot := gobot.New(
		slack.New(os.Getenv("SLACK_TOKEN")),
	)
	robot.Hear(gobot.Regexp("hi"), func(r gobot.Responder) error {
		r.Send(gobot.Message{Text: "hi to you too, " + r.User})
		return nil
	})
	robot.Run()
}

func Example_advanced() {
	robot := gobot.New(
		slack.New(os.Getenv("SLACK_TOKEN")),
	)
	robot.Enter(func(r gobot.Responder) error {
		msg := r.Message.Envelope.(*slacker.Message)

		r.Send(gobot.Message{Text: "Any friend of " + msg.Inviter + " is a friend of mine"})
		return nil
	},
	)
	robot.Run()
}

func ExampleSend() {
	adapter := slack.New(os.Getenv("SLACK_TOKEN"))
	adapter.Send(gobot.Message{Text: "hello!"})
}

func ExampleSend_custom() {
	adapter := slack.New(os.Getenv("SLACK_TOKEN"))
	adapter.Send(gobot.Message{Params: slacker.PostMessageParameters{
		Username: "ci",
		Text:     "Failed!",
		Attachments: []slacker.Attachment{
			{
				Color:     "danger",
				Title:     "CI Status",
				TitleLink: "http://ci.org/123",
				Fields: []slacker.AttachmentField{
					{Title: "Passed", Value: "102"},
					{Title: "Failed", Value: "3"},
				},
			},
		},
	}})
}

func ExampleReply() {
	adapter := slack.New(os.Getenv("SLACK_TOKEN"))
	fromMessage := gobot.Message{
		Text: "Hi bot! How are you?",
		User: "ali",
		Room: "general",
	}
	adapter.Reply(gobot.Message{
		Text:     "I'm well, thanks!",
		Envelope: fromMessage,
	})
}

func ExampleTopic() {
	adapter := slack.New(os.Getenv("SLACK_TOKEN"))
	adapter.Topic(gobot.Message{
		Room:  "general",
		Topic: "General conversation",
	})
}

func Example_Store() {
	adapter := slack.New(os.Getenv("SLACK_TOKEN"))
	// The store is only populated if:
	// 1. You call adapter.Messages(), which connects it to RTM
	adapter.Messages()
	// 2. You call store.Update()
	adapter.Store.Update()

	// Gives access to slack.User
	if u, ok := adapter.Store.UserByName("beardroid"); ok {
		fmt.Println("Found the bot's real name: " + u.RealName)
	}

	// Gives access to slack.Channel
	if c, ok := adapter.Store.ChannelByName("general"); ok {
		fmt.Println(len(c.Members), " many people in general")
	}
}
