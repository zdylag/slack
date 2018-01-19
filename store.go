package slack

type Store interface {
	UserByID(id string) (*User, bool)
	UserByName(name string) (*User, bool)
	ChannelByID(id string) (*Channel, bool)
	ChannelByName(id string) (*Channel, bool)
}

type memoryStore struct {
	cache map[string]interface{}
}

func (s *memoryStore) UserByID(id string) (*User, bool) {
	return &User{}, true
}

func (s *memoryStore) UserByName(name string) (*User, bool) {
	return &User{}, true
}

func (s *memoryStore) ChannelByID(id string) (*Channel, bool) {
	return &Channel{}, true
}

func (s *memoryStore) ChannelByName(id string) (*Channel, bool) {
	return &Channel{}, true
}
