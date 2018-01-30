package slack

import (
	"github.com/botopolis/bot"
	"github.com/nlopes/slack"
)

var onConnectNoop = func(ev *slack.ConnectedEvent) {}

type proxy struct {
	*Adapter
	RTM       *slack.RTM
	formatter formatter
	onConnect func(ev *slack.ConnectedEvent)
}

func newProxy(a *Adapter) *proxy {
	return &proxy{
		Adapter:   a,
		RTM:       a.Client.NewRTM(),
		formatter: formatter{a.Store},
		onConnect: onConnectNoop,
	}
}

func (p *proxy) OnConnect(f func(ev *slack.ConnectedEvent)) { p.onConnect = f }

func (p *proxy) Connect() chan bot.Message {
	go p.RTM.ManageConnection()
	ch := make(chan bot.Message, 32)
	go p.Forward(p.RTM.IncomingEvents, ch)
	return ch
}

func (p *proxy) Disconnect() {
	if p.RTM != nil {
		p.RTM.Disconnect()
	}
}

func (p *proxy) Forward(in <-chan slack.RTMEvent, out chan<- bot.Message) {
	defer close(out)
	for msg := range in {
		switch ev := msg.Data.(type) {
		case *slack.HelloEvent:
		case *slack.ConnectedEvent:
			p.onConnect(ev)
			p.Robot.Logger.Debugf("Connected as %s: %d", ev.Info.User.ID, ev.ConnectionCount)
		case *slack.MessageEvent:
			out <- p.translate(ev)
		case *slack.RTMError:
			p.Robot.Logger.Errorf("RTM Error: %s", ev.Error())
		case *slack.ConnectionErrorEvent:
			p.Robot.Logger.Error("Slack: Connection Error")
		case *slack.InvalidAuthEvent:
			p.Robot.Logger.Error("Slack: Invalid Credentials")
			return
		}
	}
}

func (p *proxy) translate(ev *slack.MessageEvent) bot.Message {
	user, _ := p.Store.UserByID(ev.User)
	channel, _ := p.Store.ChannelByID(ev.Channel)

	m := bot.Message{
		User:     user.Name,
		Room:     channel.Name,
		Text:     p.formatter.Format(ev),
		Topic:    ev.Topic,
		Envelope: slack.Message(*ev),
	}

	switch ev.SubType {
	case "channel_join":
		m.Type = bot.Enter
	case "channel_leave":
		m.Type = bot.Leave
	case "channel_topic":
		m.Type = bot.Topic
	default:
		m.Type = bot.DefaultMessage
	}

	return m
}
