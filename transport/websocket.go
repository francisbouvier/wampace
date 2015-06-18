package transport

import (
	log "github.com/Sirupsen/logrus"
	"github.com/francisbouvier/wampace/serializer"
	"github.com/gorilla/websocket"
)

type wsConn struct {
	serial serializer.Serializer
	conn   *websocket.Conn
}

func (wsc wsConn) Read() (msg []interface{}, err error) {
	_, data, err := wsc.conn.ReadMessage()
	if err != nil {
		return
	}
	msg, err = wsc.serial.Unmarshal(data)
	log.Debugln("Receive:", msg)
	return
}

func (wsc wsConn) Write(msg []interface{}) {
	log.Debugln("Send:", msg)
	data, err := wsc.serial.Marshal(msg)
	if wsc.conn.WriteMessage(1, data); err != nil {
		log.Errorln("Error send:", msg)
	}
}
