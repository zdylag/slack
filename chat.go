package slack

import (
	"github.com/berfarah/gobot"
	"github.com/nlopes/slack"
)

func (a *Adapter) Send(m gobot.Message) error {
	if m.Extra == nil {
		a.proxy.RTM.SendMessage(a.proxy.RTM.NewOutgoingMessage(m.Text, m.Room))
		return nil
	}

	if pm, ok := m.Extra.(slack.PostMessageParameters); ok {
		pm.AsUser = true
		if pm.User == "" {
			pm.User = a.ID
		}
		_, _, err := a.Client.PostMessage(m.Room, m.Text, pm)
		return err
	}

	return nil
}

func (a *Adapter) Reply(m gobot.Message) error {
	a.Store.UserByName(m.User)
	return nil
}
