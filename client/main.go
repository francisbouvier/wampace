package client

import (
	"crypto/tls"

	log "github.com/Sirupsen/logrus"
	"github.com/francisbouvier/wampace/serializer"
	"github.com/francisbouvier/wampace/transport"
	"github.com/francisbouvier/wampace/wamp"
)

var (
	t = "websocket"
	s = "json"
)

func get() (client transport.Client) {
	// Serializer
	serial := serializer.Json{}

	// Transport
	log.Infoln("Transport:", t)
	switch t {
	case "websocket":
		client = transport.NewWsClient(serial)
	}

	log.Infoln("Serializer:", s)
	return
}

func New(url string) (wampClient *wamp.Client, err error) {
	client := get()

	log.Infoln("Starting WAMP client on:", "ws://"+url)
	wampClient, err = client.Connect("ws://" + url)
	return
}

func NewTLS(url, certFile, keyFile string) (wampClient *wamp.Client, err error) {
	client := get()

	config := &tls.Config{
		InsecureSkipVerify: true,
	}

	log.Infoln("Starting WAMP client on:", "wss://"+url)
	wampClient, err = client.ConnectTLS("wss://"+url, config)
	return
}
