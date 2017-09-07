package http_to_nsq

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// Publisher.
type publisher interface {
	Publish(topic string, body []byte) error
}

// Message published to NSQD.
type Message struct {
	URL    string          `json:"url"`
	Method string          `json:"method"`
	Header http.Header     `json:"header"`
	Body   json.RawMessage `json:"body"`
}

// Server publishing requests as Messages.
type Server struct {
	Topic     string
	Secret    string
	Publisher publisher
	Log       *log.Logger
}

// ServeHTTP implements http.Handler.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/internal/health":
		fmt.Fprintf(w, "OK\n")
	default:
		s.publish(w, r)
	}
}

// Publish requests to NSQD.
func (s *Server) publish(w http.ResponseWriter, r *http.Request) {
	secret := r.URL.Query().Get("secret")
	if s.Secret != "" && s.Secret != secret {
		s.Log.Printf("[error] invalid secret")
		http.Error(w, http.StatusText(403), 403)
		return
	}

	var body json.RawMessage
	if r.Body != nil {
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			s.Log.Printf("[error] decoding body: %s", err)
			http.Error(w, "Error parsing request body", 400)
			return
		}
	}

	msg := &Message{
		URL:    r.URL.String(),
		Method: r.Method,
		Header: r.Header,
		Body:   body,
	}

	b, err := json.Marshal(msg)
	if err != nil {
		s.Log.Printf("[error] marshalling message: %s", err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	s.Log.Printf("[info] publishing %s %s", msg.Method, msg.URL)
	err = s.Publisher.Publish(s.Topic, b)
	if err != nil {
		s.Log.Printf("[error] publishing body: %s", err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	fmt.Fprintf(w, ":)")
}
