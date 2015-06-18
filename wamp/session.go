package wamp

import (
	"errors"
	"math/rand"
	"sync"

	log "github.com/Sirupsen/logrus"
)

var realms = map[string]map[int]*Session{}
var realmsLock sync.Mutex

type Conn interface {
	Read() ([]interface{}, error)
	Write([]interface{})
}

type Session struct {
	id     int
	status int
	realm  string
	Conn   Conn
}

func (se *Session) open(realm string) {
	se.realm = realm
	se.id = rand.Intn(RAND_MAX)
	realmsLock.Lock()
	if _, prs := realms[realm]; prs == false {
		realms[realm] = map[int]*Session{}
	}
	realmsLock.Unlock()
	realms[realm][se.id] = se
	se.status = 1
	log.Infoln("Connect peer:", se.id)
}

func (se *Session) opened(id int) {
	se.id = id
	se.status = 1
	log.Infoln("Client connected:", se.id)
}

func (se *Session) check() (err error) {
	if se.status != 1 {
		err = errors.New("Peer is not connected")
	}
	return
}

func (se *Session) Close() {
	se.status = 2
	if _, prs := realms[se.realm][se.id]; prs == true {
		delete(realms[se.realm], se.id)
	}
	// Remove peer from pubsub subscribers
	for _, ps := range psTopic {
		ps.removeSub(se)
	}
	// Remove peer from rpc callee
	for _, rp := range rps {
		rp.removeCallee(se)
	}
	log.Infoln("Disconnect peer:", se.id)
}

func (se *Session) Closed() {
	se.status = 2
	log.Infoln("Client disconnected:", se.id)
}
