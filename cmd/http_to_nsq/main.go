package main

import "github.com/segmentio/http_to_nsq"
import "github.com/bitly/go-nsq"
import "github.com/tj/docopt"
import "net/http"
import "log"
import "os"

var Version = "0.1.0"

const Usage = `
  Usage:
    http_to_nsq --topic name
      [--nsqd-tcp-address addr]
      [--address addr]
      [--secret secret]
    http_to_nsq -h | --help
    http_to_nsq --version

  Options:
    --topic name             nsqd topic name
    --secret secret          secret string [default: ]
    --address addr           bind address [default: localhost:3000]
    --nsqd-tcp-address addr  nsqd tcp address [default: localhost:4150]
    -h, --help               output help information
    -v, --version            output version

`

func main() {
	args, err := docopt.Parse(Usage, nil, true, Version, false)
	if err != nil {
		log.Fatal("[error] %s", err)
	}

	topic := args["--topic"].(string)
	addr := args["--address"].(string)
	nsqd := args["--nsqd-tcp-address"].(string)

	log.Printf("starting http_to_nsq %s", Version)
	log.Printf("--> binding to %s", addr)
	log.Printf("--> publishing to %s as topic %q", nsqd, topic)

	prod, err := nsq.NewProducer(nsqd, nsq.NewConfig())
	if err != nil {
		log.Fatal("[error] starting producer: %s", err)
	}

	server := &http_to_nsq.Server{
		Log:       log.New(os.Stderr, "", log.LstdFlags),
		Secret:    args["--secret"].(string),
		Topic:     "builds",
		Publisher: prod,
	}

	err = http.ListenAndServe(addr, server)
	if err != nil {
		log.Fatal("[error] binding: %s", err)
	}
}
