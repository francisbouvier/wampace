package wamp

import (
	"errors"

	log "github.com/Sirupsen/logrus"
)

type Client struct {
	se        *Session
	done      chan bool
	msgList   chan []interface{}
	subsId    map[int]*sub
	subs      map[int]*sub
	unsubs    map[int]*sub
	pubs      map[int]*pub
	procsId   map[int]*rProc
	procs     map[int]*rProc
	unprocs   map[int]*rProc
	calls     map[int]*RCall
	OnConnect func()
	OnJoin    func()
}

func (c *Client) do(msg []interface{}) {
	code := int(msg[0].(float64))
	switch code {
	case Subscribed:
		c.subscribed(msg)
	case Unsubscribed:
		c.unsubscribed(msg)
	case Published:
		c.published(msg)
	case Event:
		c.event(msg)
	case Registred:
		c.registred(msg)
	case Unregistred:
		c.unregistred(msg)
	case Invocation:
		c.invoke(msg)
	case Result:
		c.result(msg)
	}
}

func (c *Client) Join(realm string) (err error) {
	details := map[string]map[string]map[string]interface{}{
		"roles": {
			"publisher":  {},
			"subscriber": {},
			"caller":     {},
			"callee":     {},
		},
	}
	c.se.Conn.Write([]interface{}{Hello, realm, details})
	msg := <-c.msgList
	// Check that router answers Welcome
	code := int(msg[0].(float64))
	if code == Welcome {
		id := int(msg[1].(float64))
		c.se.opened(id)
		log.Debugln("WAMP success")
		// Launching the receive messages loop
		go c.receive()
		// Calling callback OnJoin
		c.OnJoin()
	} else {
		err = errors.New("Join error")
		// TODO: close the Wamp client and the WS connection
	}
	return
}

func (c *Client) receive() {
	log.Debugln("Lauching receive ...")
	for {
		msg := <-c.msgList
		if msg[0] == "closed" {
			c.se.Closed()
		} else {
			go c.do(msg)
		}
		if c.se.status == 2 {
			break
		}
	}
	c.done <- true
}

func (c *Client) End() {
	<-c.done
}

func NewClient(se *Session, msgList chan []interface{}) (c *Client) {
	c = &Client{
		se:        se,
		done:      make(chan bool, 1),
		msgList:   msgList,
		subsId:    map[int]*sub{},
		subs:      map[int]*sub{},
		unsubs:    map[int]*sub{},
		pubs:      map[int]*pub{},
		procsId:   map[int]*rProc{},
		procs:     map[int]*rProc{},
		unprocs:   map[int]*rProc{},
		calls:     map[int]*RCall{},
		OnConnect: func() {},
		OnJoin:    func() {},
	}
	return
}
