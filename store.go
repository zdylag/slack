package slack

import "github.com/nlopes/slack"

// Store is the interface to expect from adapter.Store
type Store interface {
	// Load takes slack info and adds new users and channels from it
	Load(*slack.Info)
	// Update queries Slack's web API for users and channels
	Update() error
	// UserByID queries the store for a User by ID
	UserByID(id string) (slack.User, bool)
	// UserByName queries the store for a User by Name
	UserByName(name string) (slack.User, bool)
	// ChannelByID queries the store for a Channel by ID
	ChannelByID(id string) (slack.Channel, bool)
	// ChannelByName queries the store for a Channel by Name
	ChannelByName(id string) (slack.Channel, bool)
}

type memoryStore struct {
	client   *slack.Client
	indices  map[string]string
	users    map[string]slack.User
	channels map[string]slack.Channel
}

func newMemoryStore(c *slack.Client) *memoryStore {
	m := &memoryStore{
		client:   c,
		indices:  make(map[string]string),
		users:    make(map[string]slack.User),
		channels: make(map[string]slack.Channel),
	}
	return m
}

func (s *memoryStore) Load(i *slack.Info) {
	for _, u := range i.Users {
		s.users[u.ID] = u
		s.indices["user:name:"+u.Name] = u.ID
	}

	for _, ch := range i.Channels {
		s.channels[ch.ID] = ch
		s.indices["channel:name:"+ch.Name] = ch.ID
	}
}

func (s *memoryStore) Update() (err error) {
	info := slack.Info{}
	if info.Users, err = s.client.GetUsers(); err != nil {
		return err
	}

	if info.Channels, err = s.client.GetChannels(true); err != nil {
		return err
	}
	s.Load(&info)
	return err
}

func (s *memoryStore) UserByID(id string) (slack.User, bool) {
	u, ok := s.users[id]
	return u, ok
}

func (s *memoryStore) UserByName(name string) (slack.User, bool) {
	return s.UserByID(s.indices["user:name:"+name])
}

func (s *memoryStore) ChannelByID(id string) (slack.Channel, bool) {
	ch, ok := s.channels[id]
	return ch, ok
}

func (s *memoryStore) ChannelByName(name string) (slack.Channel, bool) {
	return s.ChannelByID(s.indices["channel:name:"+name])
}
