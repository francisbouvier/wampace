package transport

import log "github.com/Sirupsen/logrus"

type TCP struct {
	network string
}

func (t TCP) Serve(conf *Config) {
	log.Debugln("ok", t.network)
}

func NewTCP() (t TCP) {
	t = TCP{network: "tcp"}
	return
}

func NewUnixDomain() (t TCP) {
	t = TCP{network: "unix"}
	return
}
