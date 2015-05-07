package http_to_nsq

import "github.com/bmizerany/assert"
import "net/http/httptest"
import "io/ioutil"
import "net/http"
import "testing"
import "bytes"
import "log"

type pub struct {
	msgs [][]byte
}

func (p *pub) Publish(topic string, body []byte) error {
	p.msgs = append(p.msgs, body)
	return nil
}

func TestServer_ServeHTTP_POST(t *testing.T) {
	p := new(pub)

	s := Server{
		Log:       log.New(ioutil.Discard, "", log.LstdFlags),
		Topic:     "builds",
		Publisher: p,
	}

	b := bytes.NewBufferString(`{ "foo": "bar" }`)

	r, err := http.NewRequest("POST", "/build", b)
	assert.Equal(t, nil, err)
	r.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	s.ServeHTTP(w, r)

	assert.Equal(t, 1, len(p.msgs))
	assert.Equal(t, `{"url":"/build","method":"POST","header":{"Content-Type":["application/json"]},"body":{"foo":"bar"}}`, string(p.msgs[0]))

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, ":)", w.Body.String())
}

func TestServer_ServeHTTP_malformedJSON(t *testing.T) {
	p := new(pub)

	s := Server{
		Log:       log.New(ioutil.Discard, "", log.LstdFlags),
		Topic:     "builds",
		Publisher: p,
	}

	b := bytes.NewBufferString(`{ "`)

	r, err := http.NewRequest("POST", "/build", b)
	assert.Equal(t, nil, err)
	r.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	s.ServeHTTP(w, r)

	assert.Equal(t, 0, len(p.msgs))
	assert.Equal(t, 400, w.Code)
	assert.Equal(t, "Error parsing request body\n", w.Body.String())
}
