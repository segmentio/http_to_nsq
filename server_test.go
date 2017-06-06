package http_to_nsq

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bmizerany/assert"
)

type pub struct {
	msgs [][]byte
}

func (p *pub) Publish(topic string, body []byte) error {
	p.msgs = append(p.msgs, body)
	return nil
}

func (p *pub) Reset() {
	p.msgs = nil
}

func TestServer_ServeHTTP_POST(t *testing.T) {
	p := new(pub)

	s := Server{
		Log:       log.New(ioutil.Discard, "", log.LstdFlags),
		Topic:     "builds",
		Publisher: p,
	}

	for _, method := range []string{"POST", "PUT", "DELETE"} {
		t.Run(method, func(t *testing.T) {
			defer p.Reset()

			b := bytes.NewBufferString(`{ "foo": "bar" }`)

			r, err := http.NewRequest(method, "/build", b)
			assert.Equal(t, nil, err)
			r.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			s.ServeHTTP(w, r)

			body := fmt.Sprintf(`{"url":"/build","method":"%s","header":{"Content-Type":["application/json"]},"body":{"foo":"bar"}}`, method)

			assert.Equal(t, 1, len(p.msgs))
			assert.Equal(t, body, string(p.msgs[0]))

			assert.Equal(t, 200, w.Code)
			assert.Equal(t, ":)", w.Body.String())
		})
	}
}

func TestServer_ServeHTTP_InvalidMethods(t *testing.T) {
	p := new(pub)

	s := Server{
		Log:       log.New(ioutil.Discard, "", log.LstdFlags),
		Topic:     "builds",
		Publisher: p,
	}

	for _, method := range []string{"GET", "OPTIONS", "HEAD", "CONNECT"} {
		t.Run(method, func(t *testing.T) {
			b := bytes.NewBufferString(`{ "foo": "bar" }`)
			r, err := http.NewRequest(method, "/build", b)
			assert.Equal(t, nil, err)
			r.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			s.ServeHTTP(w, r)

			assert.Equal(t, 0, len(p.msgs))

			assert.Equal(t, 405, w.Code)
			assert.Equal(t, "Method Not Allowed\n", w.Body.String())
		})
	}
}

func TestServer_ServeHTTP_secret_invalid(t *testing.T) {
	p := new(pub)

	s := Server{
		Log:       log.New(ioutil.Discard, "", log.LstdFlags),
		Topic:     "builds",
		Secret:    "wahoo",
		Publisher: p,
	}

	b := bytes.NewBufferString(`{ "foo": "bar" }`)

	r, err := http.NewRequest("POST", "/build", b)
	assert.Equal(t, nil, err)
	r.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	s.ServeHTTP(w, r)

	assert.Equal(t, 0, len(p.msgs))
	assert.Equal(t, 403, w.Code)
	assert.Equal(t, "Forbidden\n", w.Body.String())
}

func TestServer_ServeHTTP_secret_correct(t *testing.T) {
	p := new(pub)

	s := Server{
		Log:       log.New(ioutil.Discard, "", log.LstdFlags),
		Topic:     "builds",
		Secret:    "wahoo",
		Publisher: p,
	}

	b := bytes.NewBufferString(`{ "foo": "bar" }`)

	r, err := http.NewRequest("POST", "/build?secret=wahoo", b)
	assert.Equal(t, nil, err)
	r.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	s.ServeHTTP(w, r)

	assert.Equal(t, 1, len(p.msgs))
	assert.Equal(t, 200, w.Code)
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

func TestServer_ServeHTTP_health(t *testing.T) {
	p := new(pub)

	s := Server{
		Log:       log.New(ioutil.Discard, "", log.LstdFlags),
		Topic:     "builds",
		Publisher: p,
	}

	r, err := http.NewRequest("GET", "/internal/health", nil)
	assert.Equal(t, nil, err)

	w := httptest.NewRecorder()
	s.ServeHTTP(w, r)

	assert.Equal(t, 0, len(p.msgs))
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "OK\n", w.Body.String())
}

type responseWriterNoOp struct {
}

func (*responseWriterNoOp) Header() http.Header {
	return nil
}

func (*responseWriterNoOp) Write([]byte) (int, error) {
	return 0, nil
}

func (*responseWriterNoOp) WriteHeader(int) {
}

func BenchmarkServe(b *testing.B) {
	p := new(pub)

	s := Server{
		Log:       log.New(ioutil.Discard, "", log.LstdFlags),
		Topic:     "builds",
		Publisher: p,
	}

	for i := 0; i < b.N; i++ {
		r, err := http.NewRequest("POST", "/build", bytes.NewBufferString(`{ "foo": "bar" }`))
		if err != nil {
			b.Error(err)
		}
		r.Header.Set("Content-Type", "application/json")
		s.ServeHTTP(&responseWriterNoOp{}, r)
	}
}
