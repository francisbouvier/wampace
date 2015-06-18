package transport

import (
	"crypto/tls"

	"github.com/francisbouvier/wampace/wamp"
)

type Config struct {
	Addr    string
	Scheme  string
	TLS     bool
	TLSCert string
	TLSKey  string
}

type Transport interface {
	Serve(*Config)
}

type Client interface {
	Connect(string) (*wamp.Client, error)
	ConnectTLS(string, *tls.Config) (*wamp.Client, error)
}
