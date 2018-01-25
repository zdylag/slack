package slack

import (
	"errors"
	"sync"

	"github.com/berfarah/gobot"
	"github.com/nlopes/slack"
)

// Adapter is the gobot slack adapter it implements
// gobot.Plugin and gobot.Chat interfaces
type Adapter struct {
	mu    sync.Mutex
	proxy *proxy

	Robot  *gobot.Robot
	Client *slack.Client
	Store  Store

	ID    string
	BotID string
	Name  string
}

// New called with one's slack token provides a new adapter
func New(secret string) *Adapter {
	client := slack.New(secret)
	a := &Adapter{
		Client: client,
		Store:  newMemoryStore(client),
	}
	a.proxy = newProxy(a)
	a.proxy.OnConnect(func(ev *slack.ConnectedEvent) {
		a.mu.Lock()
		defer a.mu.Unlock()
		a.Store.Load(ev.Info)
		u := ev.Info.User
		a.BotID = u.ID
		a.Name = u.Name
	})
	return a
}

// Load provides the slack adapter access to the Robot's logger
func (a *Adapter) Load(r *gobot.Robot) { a.Robot = r }

// Username returns the bot's username
func (a *Adapter) Username() string { return a.Name }

// Messages connects to Slack's RTM API and channels messages through
func (a *Adapter) Messages() <-chan gobot.Message { return a.proxy.Connect() }

func emptyMessage(m gobot.Message) bool {
	return m.Text == "" && m.Params == nil
}

// Send send messages to Slack. If only text is provided, it uses
// the already open RTM connection. If slack.PostMessageParamters
// are provided in the message.Params field, it will send a web
// API request.
func (a *Adapter) Send(m gobot.Message) error {
	if emptyMessage(m) {
		return nil
	}

	if m.Params == nil {
		a.proxy.RTM.SendMessage(a.proxy.RTM.NewOutgoingMessage(m.Text, m.Room))
		return nil
	}

	if pm, ok := m.Params.(slack.PostMessageParameters); ok {
		pm.AsUser = true
		if pm.User == "" {
			pm.User = a.ID
		}
		_, _, err := a.Client.PostMessage(m.Room, m.Text, pm)
		return err
	}

	return nil
}

// Reply does the same thing as send, but prefixes the message
// with <@userID>, notifying the user of the message.
func (a *Adapter) Reply(m gobot.Message) error {
	if emptyMessage(m) {
		return nil
	}

	msg, ok := m.Envelope.(slack.Message)
	if !ok || len(msg.Channel) == 0 {
		return errors.New("No Envelope provided")
	}

	// No need to @ the user if it's a DM
	if msg.Channel[0] != 'D' {
		m.Text = "<@" + msg.User + ">" + m.Text
	}

	return a.Send(m)
}

// Topic uses the web API to change the topic. It prefers
// the message.Room and falls back to message.Extra.Channel
// to determine what channel's topic should be updated.
func (a *Adapter) Topic(m gobot.Message) error {
	var channelID string

	if m.Room == "" {
		msg, ok := m.Envelope.(slack.Message)
		if !ok {
			return errors.New("No Channel provided")
		}
		channelID = msg.Channel
	} else {
		ch, ok := a.Store.ChannelByName(m.Room)
		if !ok {
			return errors.New("Channel not found")
		}
		channelID = ch.ID
	}
	_, err := a.Client.SetChannelTopic(channelID, m.Topic)

	return err
}
