package slack

import (
	"github.com/berfarah/gobot"
	"github.com/nlopes/slack"
)

type Adapter struct {
	proxy *proxy

	Robot  *gobot.Robot
	Client *slack.Client
	Store  Store

	ID    string
	BotID string
	Name  string
}

func New(r *gobot.Robot, secret string) *Adapter {
	a := &Adapter{
		Robot:  r,
		Client: slack.New(secret),
		Store:  &memoryStore{},
	}
	a.proxy = newProxy(a)
	a.proxy.OnConnect(func(ev *slack.ConnectedEvent) {
		u := ev.Info.User
		a.BotID = u.ID
		a.Name = u.Name
		r.Name = u.Name
	})
	return a
}

func (a *Adapter) Messages() chan *gobot.MessageEvent {
	return a.proxy.Connect()
}
