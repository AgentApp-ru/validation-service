package main

import (
	"net"

	"github.com/bshuster-repo/logrus-logstash-hook"
	"github.com/sirupsen/logrus"
)

func main() {
	log := logrus.New()
	conn, err := net.Dial("tcp", "elk.b2bpolis.ru:5000")
	if err != nil {
		log.Fatal(err)
	}
	hook := logrustash.New(conn, logrustash.DefaultFormatter(logrus.Fields{"type": "validation"}))

	log.Hooks.Add(hook)
	ctx := log.WithFields(logrus.Fields{
		"partner_store": "lo-vf",
	})
	ctx.Info("Hello World!")
}
