# Wampace

Wampace is [WAMP](http://wamp.ws/) router and client written in Go.

## Installation

### From source

```sh
go get github.com/francisbouvier/wampace
```

## Router usage

```sh
wampace
# Will launch the router on 127.0.0.1:1234
```

**Address and port**

You can specify on witch address/port you want to use the router.

```sh
wampace <ip>:<port>
```

**TLS Mode**

In ordrer to use the router on TLS (SSL) mode, you need to specify the path to the certificate and key.
Paths can be relative or absolute.

```sh
wampace -c server.crt -k key.pem
```

*NB*: You can use either self-signed certificates or commercial certificates. Self-signed certificates will not work with a WAMP browser clients (except if you add them manually in the brower).

## Client usage

```go
package main

import (
	"log"

	"github.com/francisbouvier/wampace/client"
)

func main() {

	// Set the address of the router
	addr := "127.0.0.1:1234"

	// Define a realm
	realm := "realm1"

	// Connect to the server
	c, err := client.New(addr)
	if err != nil {
		log.Fatalln("Connect error:", err)
	}

	// Join the wamp session
	err = c.Join(realm)
	if err != nil {
		log.Fatalln("Join error:", err)
	}

	// Wait for messages
	c.End()
}
```

See the `examples` directory for detailled pubsub and rpc examples.
