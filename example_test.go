package slack_test

import (
	"fmt"
	"os"

	"github.com/botopolis/bot"
	"github.com/botopolis/slack"
	slacker "github.com/nlopes/slack"
)

func Example() {
	robot := bot.New(
		slack.New(os.Getenv("SLACK_TOKEN")),
	)
	robot.Enter(func(r bot.Responder) error {
		msg := r.Message.Envelope.(*slacker.Message)

		r.Send(bot.Message{Text: "Any friend of " + msg.Inviter + " is a friend of mine"})
		return nil
	},
	)
	robot.Run()
}

func ExampleSend() {
	adapter := slack.New(os.Getenv("SLACK_TOKEN"))
	adapter.Send(bot.Message{Text: "hello!"})
}

func ExampleSend_custom() {
	adapter := slack.New(os.Getenv("SLACK_TOKEN"))
	adapter.Send(bot.Message{Params: slacker.PostMessageParameters{
		Username: "ci",
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
	fromMessage := bot.Message{
		Text: "Hi bot! How are you?",
		User: "ali",
		Room: "general",
	}
	adapter.Reply(bot.Message{
		Text:     "I'm well, thanks!",
		Envelope: fromMessage,
	})
}

func ExampleTopic() {
	adapter := slack.New(os.Getenv("SLACK_TOKEN"))
	adapter.Topic(bot.Message{
		Room:  "general",
		Topic: "General conversation",
	})
}

func ExampleStore() {
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
