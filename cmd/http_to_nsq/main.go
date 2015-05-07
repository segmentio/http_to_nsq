package main

import "github.com/segmentio/http_to_nsq"
import "github.com/bitly/go-nsq"
import "github.com/tj/docopt"

// import "net/http"
import "log"

var Version = "0.0.1"

const Usage = `
  Usage:
    http_to_nsq [--nsqd-tcp-address addr] [--address addr]
    http_to_nsq -h | --help
    http_to_nsq --version

  Options:
    --address addr           bind address [default: localhost:3000]
    --nsqd-tcp-address addr  nsqd tcp address [default: localhost:4150]
    -h, --help               output help information
    -v, --version            output version

`

func main() {
	_, err := docopt.Parse(Usage, nil, true, Version, false)
	if err != nil {
		log.Fatal("error: %s", err)
	}

	// nsqd := args["--nsqd-tcp-address"].(string)
	// prod, err := nsq.NewProducer(nsqd, nsq.NewConfig())
	// if err != nil {
	// 	log.Fatal("error starting producer: %s", err)
	// }

	// err := http.ListenAndServe(addr, handler)
	// if err != nil {
	// 	log.Fatal("error binding: %s", err)
	// }
}
