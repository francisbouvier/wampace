package wamp

import (
	"math/rand"

	log "github.com/Sirupsen/logrus"
)

func (c *Client) Register(proc string, cbk proc) (p *rProc) {
	p = &rProc{
		regRequest:  rand.Intn(RAND_MAX),
		procedure:   proc,
		cbk:         cbk,
		Registred:   make(chan error),
		Unregistred: make(chan error),
	}
	c.procs[p.regRequest] = p
	msg := []interface{}{Register, p.regRequest, details{}, proc}
	c.se.Conn.Write(msg)
	return
}

func (c *Client) Unregister(p *rProc) {
	p.unregRequest = rand.Intn(RAND_MAX)
	c.unprocs[p.unregRequest] = p
	msg := []interface{}{Unregistred, p.unregRequest, p.id}
	c.se.Conn.Write(msg)
}

func (c *Client) registred(msg []interface{}) {
	regRequest := int(msg[1].(float64))
	p := c.procs[regRequest]
	p.id = int(msg[2].(float64))
	c.procsId[p.id] = p
	var err error
	p.Registred <- err
	log.Debugln("Registration ok")
}

func (c *Client) unregistred(msg []interface{}) {
	unregRequest := int(msg[1].(float64))
	p := c.unprocs[unregRequest]
	delete(c.procsId, p.id)
	delete(c.procs, p.regRequest)
	delete(c.unprocs, p.unregRequest)
	var err error
	p.Unregistred <- err
	log.Debugln("Unregistration ok")
}

func (c *Client) Call(proc string, a args, k kwargs) (r *RCall) {
	r = &RCall{
		request: rand.Intn(RAND_MAX),
		Result:  make(chan error),
	}
	c.calls[r.request] = r
	msg := []interface{}{Call, r.request, details{}, proc, a, k}
	c.se.Conn.Write(msg)
	return
}

func (c *Client) invoke(msg []interface{}) {
	id := int(msg[2].(float64))
	p := c.procsId[id]
	var a args
	var k kwargs
	if len(msg) > 3 {
		a = msg[4].([]interface{})
	}
	if len(msg) > 4 {
		k = msg[5].(map[string]interface{})
	}
	ar, kr := p.cbk(a, k)
	// Yield result
	request := int(msg[1].(float64))
	m := []interface{}{Yield, request, details{}, ar, kr}
	c.se.Conn.Write(m)
	log.Debugln("Invocation ok")
}

func (c *Client) result(msg []interface{}) {
	request := int(msg[1].(float64))
	r := c.calls[request]
	delete(c.calls, r.request)
	if len(msg) > 3 {
		r.Args = msg[3].([]interface{})
	}
	// TODO: Kwargs, router pb ??
	// if len(msg) > 4 {
	// 	r.Kwargs = msg[4].(map[string]interface{})
	// }
	var err error
	r.Result <- err
	return
}
