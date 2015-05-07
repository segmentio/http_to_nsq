package http_to_nsq

import "encoding/json"
import "net/http"
import "fmt"
import "log"

// Publisher.
type publisher interface {
	Publish(topic string, body []byte) error
}

// Message published to NSQD.
type Message struct {
	URL    string                 `json:"url"`
	Method string                 `json:"method"`
	Header http.Header            `json:"header"`
	Body   map[string]interface{} `json:"body"`
}

// Server publishing requests as Messages.
type Server struct {
	Topic     string
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

func (s *Server) publish(w http.ResponseWriter, r *http.Request) {
	var body map[string]interface{}

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		s.Log.Printf("error decoding body: %s", err)
		http.Error(w, "Error parsing request body", 400)
		return
	}

	msg := &Message{
		URL:    r.URL.String(),
		Method: r.Method,
		Header: r.Header,
		Body:   body,
	}

	b, err := json.Marshal(msg)
	if err != nil {
		s.Log.Printf("error marshalling message: %s", err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	err = s.Publisher.Publish(s.Topic, b)
	if err != nil {
		s.Log.Printf("error publishing body: %s", err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	fmt.Fprintf(w, ":)")
}
