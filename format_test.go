package slack // don't want to export, do want to test
import (
	"testing"

	"github.com/nlopes/slack"
	"github.com/stretchr/testify/assert"
)

type formatterStore struct{}

func (f formatterStore) UserByID(s string) (*User, bool) {
	return &User{ID: s, Name: "bob"}, true
}
func (f formatterStore) UserByName(s string) (*User, bool) {
	return &User{ID: "U1234", Name: "bob"}, true
}
func (f formatterStore) ChannelByID(s string) (*Channel, bool) {
	c := &Channel{}
	c.ID = s
	c.Name = "general"
	return c, true
}
func (f formatterStore) ChannelByName(s string) (*Channel, bool) {
	c := &Channel{}
	c.ID = "C1234"
	c.Name = s
	return c, true
}

func TestFormatter(t *testing.T) {
	assert := assert.New(t)

	f := formatter{formatterStore{}}

	assert.Equal("", f.Format(&slack.MessageEvent{}))
}
