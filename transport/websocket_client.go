package transport

import (
	"crypto/tls"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/francisbouvier/wampace/serializer"
	"github.com/francisbouvier/wampace/wamp"
	"github.com/gorilla/websocket"
)

// wsClientConn

type wsClientConn struct {
	serial  serializer.Serializer
	msgList chan []interface{}
	se      *wamp.Session
	conn    *websocket.Conn
}

func (wsc *wsClientConn) connect(url string, dialer *websocket.Dialer) (c *wamp.Client, err error) {
	h := http.Header{}
	wsc.conn, _, err = dialer.Dial(url, h)
	if err != nil {
		return
	}
	log.Debugln("WebSocket success")
	wsc.se = &wamp.Session{
		Conn: wsConn{
			conn:   wsc.conn,
			serial: wsc.serial,
		},
	}
	go wsc.catch()
	c = wamp.NewClient(wsc.se, wsc.msgList)
	c.OnConnect()
	return
}

func (wsc *wsClientConn) catch() {
	log.Debugln("Lauching stream ...")
	for {
		msg, err := wsc.se.Conn.Read()
		if err != nil {
			log.Debugln("Close connection")
			wsc.msgList <- []interface{}{"closed"}
			return
		}
		wsc.msgList <- msg
	}
}

// wsClient

type wsClient struct {
	wsc *wsClientConn
}

func (wsci wsClient) Connect(url string) (c *wamp.Client, err error) {
	dialer := websocket.DefaultDialer
	return wsci.wsc.connect(url, dialer)
}

func (wsci wsClient) ConnectTLS(url string, config *tls.Config) (c *wamp.Client, err error) {
	dialer := &websocket.Dialer{TLSClientConfig: config}
	return wsci.wsc.connect(url, dialer)
}

func NewWsClient(serial serializer.Serializer) (wsc wsClient) {
	wsc = wsClient{
		wsc: &wsClientConn{
			serial:  serial,
			msgList: make(chan []interface{}, 1),
		},
	}
	return
}
