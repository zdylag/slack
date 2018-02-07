package action

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nlopes/slack"
	"github.com/stretchr/testify/assert"
)

func readCloser(b []byte) io.ReadCloser {
	return ioutil.NopCloser(bytes.NewReader(b))
}

func TestWebhook_response(t *testing.T) {
	token := "asjdkflj120fj0912"
	cases := []struct {
		In  []byte
		Out int
	}{
		{
			In:  []byte(`<xml></xml>`),
			Out: http.StatusBadRequest,
		},
		{
			In:  []byte(`{}`),
			Out: http.StatusBadRequest,
		},
		{
			In:  []byte(`{"token":"foo"}`),
			Out: http.StatusBadRequest,
		},
		{
			In:  []byte(`{"token":"` + token + `"}`),
			Out: http.StatusOK,
		},
	}

	p := Plugin{Token: token}
	for _, c := range cases {
		recorder := httptest.NewRecorder()
		p.webhook(recorder, &http.Request{Body: readCloser(c.In)})
		assert.Equal(t, c.Out, recorder.Code)
	}
}

func TestWebhook_callback(t *testing.T) {
	token := "asjdkflj120fj0912"
	done := make(chan string)
	fooReq := http.Request{Body: readCloser(
		[]byte(`{"callback_id":"foo", "token": "` + token + `"}`),
	)}
	barReq := http.Request{Body: readCloser(
		[]byte(`{"callback_id":"bar", "token": "` + token + `"}`),
	)}

	p := Plugin{Token: token}
	p.Add("bar", func(slack.AttachmentActionCallback) { done <- "bar" })
	p.Add("foo", func(slack.AttachmentActionCallback) { done <- "foo" })

	p.webhook(httptest.NewRecorder(), &fooReq)
	assert.Equal(t, "foo", <-done)

	p.webhook(httptest.NewRecorder(), &barReq)
	assert.Equal(t, "bar", <-done)
}
