package transport

import (
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/francisbouvier/wampace/serializer"
	"github.com/francisbouvier/wampace/wamp"
	"github.com/gorilla/websocket"
)

func checkOrigin(r *http.Request) bool {
	return true
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     checkOrigin,
}

type wsHandler struct {
	serial serializer.Serializer
	server wamp.Server
}

func (wsh *wsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// Starting connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Errorln(err)
		return
	}
	se := &wamp.Session{
		Conn: wsConn{
			conn:   conn,
			serial: wsh.serial,
		},
	}

	//  Messages
	for {
		msg, err := se.Conn.Read()
		if err != nil {
			log.Debugln("Close connection")
			se.Close()
			return
		}
		resp, err := wsh.server.Route(se, msg)
		if len(resp) > 0 {
			se.Conn.Write(resp)
		}
	}
}

type wsServer struct {
	serial serializer.Serializer
}

func (ws wsServer) Serve(config *Config) {
	wsh := &wsHandler{
		serial: ws.serial,
		server: wamp.NewServer(),
	}
	http.Handle("/", wsh)
	log.Infoln("Starting WAMP server on:", config.Scheme+config.Addr)
	var err error
	if config.TLS {
		err = http.ListenAndServeTLS(config.Addr, config.TLSCert, config.TLSKey, nil)
	} else {
		err = http.ListenAndServe(config.Addr, nil)
	}
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func NewWsServer(serial serializer.Serializer) (ws wsServer) {
	ws = wsServer{
		serial: serial,
	}
	return
}
