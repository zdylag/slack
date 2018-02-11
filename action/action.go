package action

import (
	"encoding/json"
	"net/http"

	"github.com/botopolis/bot"
	"github.com/nlopes/slack"
)

// Plugin conforms to the botopolis/bot.Plugin interface
type Plugin struct {
	*registry
	// Path at which our webhook sits
	Path string
	// Token to verify message comes from slack
	Token string
}

// New returns a new plugin taking arguments for path and token
func New(path, token string) *Plugin {
	return &Plugin{
		registry: &registry{},
		Path:     path,
		Token:    token,
	}
}

// Load installs the webhook
func (p Plugin) Load(r *bot.Robot) {
	r.Router.HandleFunc(p.Path, p.webhook)
}

func (p Plugin) webhook(w http.ResponseWriter, r *http.Request) {
	var cb slack.AttachmentActionCallback
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&cb); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if cb.Token != p.Token {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	go p.Run(cb)
}
