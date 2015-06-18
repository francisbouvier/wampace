# Wampace, a WAMP router

Wampace is WAMP routeur written in Go.

## Installation

### From source

```sh
go get github.com/francisbouvier/wampace
```

## Usage

```sh
wampace
# Will launch the router on localhost:1234
```

**Address and port**

You can specify on witch address/port you want to use the router.

```sh
wampace mydomain.com
wampace mydomain.com:8888
```

**TLS Mode**

In ordrer to use the router on TLS (SSL) mode, you need to specify the path to the certificate and key.
Paths can be relative or absolute.

```sh
wampace -c server.crt -k key.pem
```

*NB*: You can use either self-signed certificates or commercial certificates. Self-signed certificates will not work with a WAMP browser clients (except if you add them manually in the brower).
