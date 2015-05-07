package http_to_nsq

import "encoding/json"
import "net/http"
import "fmt"
import "log"

type Publisher interface {
	Publish(topic string, body []byte) error
}

type Message struct {
	URL    string                 `json:"url"`
	Method string                 `json:"method"`
	Header http.Header            `json:"header"`
	Body   map[string]interface{} `json:"body"`
}

type Server struct {
	Topic     string
	Publisher Publisher
	Log       *log.Logger
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
