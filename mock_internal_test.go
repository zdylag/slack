package slack

import (
	"github.com/nlopes/slack"
)

type testStore struct {
	LoadFunc   func(*slack.Info)
	UpdateFunc func() error
	User       slack.User
	Channel    slack.Channel
	IM         slack.IM
}

func newTestStore() *testStore {
	return &testStore{
		LoadFunc:   func(*slack.Info) {},
		UpdateFunc: func() error { return nil },
		User:       slack.User{},
		Channel:    slack.Channel{},
		IM:         slack.IM{},
	}
}

func (s *testStore) Load(i *slack.Info) { s.LoadFunc(i) }
func (s *testStore) Update() error      { return s.UpdateFunc() }
func (s *testStore) UserByID(id string) (slack.User, bool) {
	if s.User.ID == id {
		return s.User, true
	}
	return s.User, false
}
func (s *testStore) UserByName(name string) (slack.User, bool) {
	if s.User.Name == name {
		return s.User, true
	}
	return s.User, false
}
func (s *testStore) ChannelByID(id string) (slack.Channel, bool) {
	if s.Channel.ID == id {
		return s.Channel, true
	}
	return s.Channel, false
}
func (s *testStore) ChannelByName(name string) (slack.Channel, bool) {
	if s.Channel.Name == name {
		return s.Channel, true
	}
	return s.Channel, false
}
func (s *testStore) IMByID(id string) (slack.IM, bool) {
	if s.IM.ID == id {
		return s.IM, true
	}
	return s.IM, false
}
func (s *testStore) IMByUserID(id string) (slack.IM, bool) {
	if s.IM.User == id {
		return s.IM, true
	}
	return s.IM, false
}
