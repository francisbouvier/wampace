package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/francisbouvier/wampace/server"
	"github.com/francisbouvier/wampace/transport"
)

var usage = `Wampace, a WAMP router.

Usage:
  wampace [options] addr
  wampace [options] addr:port

Options:`

func cliError(err error) {
	fmt.Println(err)
	flag.Usage()
	os.Exit(2)
}

func checkOption(k string, o []string, v string) {
	prs := false
	for _, c := range o {
		if c == v {
			prs = true
			break
		}
	}
	if prs == false {
		msg := fmt.Sprintf("Available choices (%s): -%s\n", strings.Join(o, ", "), k)
		err := errors.New(msg)
		cliError(err)
	}
}

func main() {

	// Help
	flag.Usage = func() {
		fmt.Println(usage)
		flag.PrintDefaults()
	}

	// Options
	d := flag.Bool("d", false, "Debug mode.")
	t := flag.String("t", "websocket", "Transport (websocket, tcp, unix).")
	s := flag.String("s", "json", "Serializer (json, msgpack).")
	c := flag.String("c", "", "TLS certificate (absolute or relative path).")
	k := flag.String("k", "", "TLS key (absolute or relative path).")
	flag.Parse()
	checkOption("s", []string{"json", "msgpack"}, *s)
	checkOption("t", []string{"websocket", "tcp", "unix"}, *t)

	// Debug mode
	if *d {
		log.SetLevel(log.DebugLevel)
	}

	// Config
	var (
		config = &transport.Config{}
		addr   string
	)
	switch *t {
	default:
		// Endpoint
		port := "1234"
		addr = "0.0.0.0"
		if len(flag.Args()) > 0 {
			addr = flag.Args()[0]
		}
		if prs := strings.Contains(addr, ":"); prs == false {
			addr = strings.Join([]string{addr, port}, ":")
		}
		config.Addr = addr

		// TLS
		if *c != "" {
			var tlsErr error
			for k, v := range map[string]string{"c": *c, "k": *k} {
				if _, err := os.Stat(v); os.IsNotExist(err) {
					if tlsErr == nil {
						tlsErr = errors.New("")
					}
					e := "File does not exist: -" + k + "\n"
					tlsErr = errors.New(tlsErr.Error() + e)
				}
			}
			if tlsErr == nil {
				log.Debugln("TLS mode")
				config.TLS = true
				config.TLSCert = *c
				config.TLSKey = *k
				config.Scheme = "wss://"
			}
			if tlsErr != nil {
				cliError(tlsErr)
			}
		} else {
			config.Scheme = "ws://"
		}
	}

	// Serve
	server := server.New(*t, *s)
	server.Serve(config)
}
