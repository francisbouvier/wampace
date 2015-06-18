package wamp

import (
	"math/rand"

	log "github.com/Sirupsen/logrus"
)

func (c *Client) Subscribe(topic string, cbk callback) (s *sub) {
	s = &sub{
		subRequest:   rand.Intn(RAND_MAX),
		topic:        topic,
		cbk:          cbk,
		Subscribed:   make(chan error),
		Unsubscribed: make(chan error),
	}
	c.subs[s.subRequest] = s
	msg := []interface{}{Subscribe, s.subRequest, details{}, topic}
	c.se.Conn.Write(msg)
	return
}

func (c *Client) Unsubscribe(s *sub) {
	s.unsubRequest = rand.Intn(RAND_MAX)
	msg := []interface{}{Unsubscribe, s.unsubRequest, s.id}
	c.unsubs[s.unsubRequest] = s
	c.se.Conn.Write(msg)
}

func (c *Client) subscribed(msg []interface{}) {
	subRequest := int(msg[1].(float64))
	s := c.subs[subRequest]
	s.id = int(msg[2].(float64))
	c.subsId[s.id] = s
	var err error
	s.Subscribed <- err
	log.Debugln("Subscription ok")
}

func (c *Client) unsubscribed(msg []interface{}) {
	unsubRequest := int(msg[1].(float64))
	s := c.unsubs[unsubRequest]
	delete(c.subsId, s.id)
	delete(c.subs, s.subRequest)
	delete(c.unsubs, s.unsubRequest)
	var err error
	s.Unsubscribed <- err
	log.Debugln("Unsubscription ok")
}

func (c *Client) Publish(topic string, a args, k kwargs) (p *pub) {
	p = &pub{
		request:   rand.Intn(RAND_MAX),
		a:         a,
		k:         k,
		Published: make(chan error),
	}
	c.pubs[p.request] = p
	msg := []interface{}{Publish, p.request, details{}, topic, a, k}
	c.se.Conn.Write(msg)
	return
}

func (c *Client) published(msg []interface{}) {
	request := int(msg[1].(float64))
	p := c.pubs[request]
	p.id = int(msg[2].(float64))
	var err error
	p.Published <- err
	log.Debugln("Publication ok")
}

func (c *Client) event(msg []interface{}) {
	subId := int(msg[1].(float64))
	s := c.subsId[subId]
	var a args
	var k kwargs
	if len(msg) > 3 {
		a = msg[4].([]interface{})
	}
	if len(msg) < 4 {
		k = msg[5].(map[string]interface{})
	}
	s.cbk(a, k)
	log.Debugln("Event ok")
}
