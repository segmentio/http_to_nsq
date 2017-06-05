package main

import (
	"log"
	"net/http"
	"os"

	"github.com/nsqio/go-nsq"
	"github.com/segmentio/http_to_nsq"
	"github.com/tj/docopt"
)

var Version = "0.2.0"

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
	secret := args["--secret"].(string)

	log.Printf("starting http_to_nsq %s", Version)
	log.Printf("--> binding to %s", addr)
	log.Printf("--> publishing to %s as topic %q", nsqd, topic)

	prod, err := nsq.NewProducer(nsqd, nsq.NewConfig())
	if err != nil {
		log.Fatal("[error] starting producer: %s", err)
	}

	server := &http_to_nsq.Server{
		Log:       log.New(os.Stderr, "", log.LstdFlags),
		Secret:    secret,
		Topic:     topic,
		Publisher: prod,
	}

	err = http.ListenAndServe(addr, server)
	if err != nil {
		log.Fatal("[error] binding: %s", err)
	}
}
