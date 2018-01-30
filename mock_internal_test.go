package slack

import "github.com/nlopes/slack"

type testStore struct{}

func (f testStore) Load(i *slack.Info) {}
func (f testStore) Update() error      { return nil }
func (f testStore) UserByID(s string) (slack.User, bool) {
	return slack.User{ID: s, Name: "bob"}, true
}
func (f testStore) UserByName(s string) (slack.User, bool) {
	return slack.User{ID: "U1234", Name: "bob"}, true
}
func (f testStore) ChannelByID(s string) (slack.Channel, bool) {
	c := slack.Channel{}
	c.ID = s
	c.Name = "general"
	return c, true
}
func (f testStore) ChannelByName(s string) (slack.Channel, bool) {
	c := slack.Channel{}
	c.ID = "C1234"
	c.Name = s
	return c, true
}
func (f testStore) IMByID(s string) (slack.IM, bool) {
	im := slack.IM{}
	im.ID = s
	im.User = "U1234"
	return im, true
}
func (f testStore) IMByUserID(s string) (slack.IM, bool) {
	im := slack.IM{}
	im.ID = "D1234"
	im.User = s
	return im, true
}
