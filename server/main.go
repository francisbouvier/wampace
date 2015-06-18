package server

import (
	log "github.com/Sirupsen/logrus"

	"github.com/francisbouvier/wampace/serializer"
	"github.com/francisbouvier/wampace/transport"
)

func New(t, s string) (trans transport.Transport) {

	// Serializer
	var serial serializer.Serializer
	serial = serializer.Json{}

	// Transport
	log.Infoln("Transport:", t)
	switch t {
	case "tcp":
		trans = transport.NewTCP()
	case "unix":
		trans = transport.NewUnixDomain()
	default:
		trans = transport.NewWsServer(serial)
	}

	// Serve
	log.Infoln("Serializer:", s)
	return
}
