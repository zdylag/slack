package slack // don't want to export, do want to test
import (
	"testing"

	"github.com/nlopes/slack"
	"github.com/stretchr/testify/assert"
)

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

func TestFormatter(t *testing.T) {
	assert := assert.New(t)

	f := formatter{testStore{}}

	assert.Equal("", f.Format(&slack.MessageEvent{}))
}
