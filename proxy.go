package slack

import (
	"github.com/berfarah/gobot"
	"github.com/nlopes/slack"
)

var onConnectNoop = func(ev *slack.ConnectedEvent) {}

type proxy struct {
	RTM       *slack.RTM
	adapter   *Adapter
	formatter formatter
	onConnect func(ev *slack.ConnectedEvent)
}

func newProxy(a *Adapter) *proxy {
	return &proxy{
		RTM:       a.Client.NewRTM(),
		adapter:   a,
		formatter: formatter{a.Store},
		onConnect: onConnectNoop,
	}
}

func (p proxy) OnConnect(f func(ev *slack.ConnectedEvent)) { p.onConnect = f }

func (p proxy) Connect() chan *gobot.MessageEvent {
	go p.RTM.ManageConnection()
	ch := make(chan *gobot.MessageEvent)
	go p.forwardEvents(ch)
	return ch
}

func (p proxy) forwardEvents(ch chan *gobot.MessageEvent) {
	defer close(ch)
	for msg := range p.RTM.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.HelloEvent:
		case *slack.ConnectedEvent:
			p.onConnect(ev)
			p.adapter.Robot.Logger.Debugf("Connected as %s: %d", ev.Info.User.ID, ev.ConnectionCount)
		case *slack.MessageEvent:
			ch <- p.translate(ev)
		case *slack.RTMError:
			p.adapter.Robot.Logger.Errorf("RTM Error: %s", ev.Error())
		case *slack.InvalidAuthEvent:
			p.adapter.Robot.Logger.Error("Slack: Invalid Credentials")
			return
		}
	}
}

func (p proxy) translate(ev *slack.MessageEvent) (me *gobot.MessageEvent) {
	user, _ := p.adapter.Store.UserByID(ev.User)
	channel, _ := p.adapter.Store.ChannelByID(ev.Channel)
	m := gobot.Message{
		User:  user.Name,
		Room:  channel.Name,
		Text:  p.formatter.Format(ev),
		Extra: ev,
	}
	me.Type = ev.SubType

	switch ev.SubType {
	case "channel_join":
		tmp := gobot.EnterMessage(m)
		me.Data = &tmp
	case "channel_leave":
		tmp := gobot.LeaveMessage(m)
		me.Data = &tmp
	case "channel_topic":
		me.Data = &gobot.TopicChange{m, ev.Topic}
	case "":
		me.Type = ev.Type
		me.Data = &m
	}

	return me
}
