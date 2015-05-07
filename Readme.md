# http_to_nsq

 Publishes HTTP requests to NSQD.

## Usage

#### type Message

```go
type Message struct {
  URL    string                 `json:"url"`
  Method string                 `json:"method"`
  Header http.Header            `json:"header"`
  Body   map[string]interface{} `json:"body"`
}
```
