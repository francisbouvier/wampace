package wamp

import (
	"math/rand"
	"sync"

	log "github.com/Sirupsen/logrus"
)

var (
	psTopic = map[string]*pub{}
	psID    = map[int]*pub{}
	psLock  sync.Mutex
)

func (s Server) subscribed(se *Session, msg []interface{}) (resp []interface{}, err error) {
	request := int(msg[1].(float64))
	topic := msg[3].(string)
	var ps pub
	psLock.Lock()
	if _, prs := psTopic[topic]; prs == false {
		ps = pub{
			id:          rand.Intn(RAND_MAX),
			topic:       topic,
			subscribers: map[int]*Session{},
		}
		psTopic[topic] = &ps
		psID[ps.id] = &ps
	} else {
		ps = *psTopic[topic]
	}
	psLock.Unlock()
	ps.addSub(se)
	resp = []interface{}{Subscribed, request, ps.id}
	return
}

func (s Server) unsubscribed(se *Session, msg []interface{}) (resp []interface{}, err error) {
	request := int(msg[1].(float64))
	id := int(msg[2].(float64))
	ps, prs := psID[id]
	if prs == true {
		ps.removeSub(se)
		resp = []interface{}{Unsubscribed, request}
	} else {
		details := map[string]string{}
		uri := "wamp.error.no_such_subscription"
		resp = []interface{}{Error, Unsubscribe, request, details, uri}
	}
	return
}

func (s Server) published(se *Session, msg []interface{}) (resp []interface{}, err error) {
	request := int(msg[1].(float64))
	topic := msg[3].(string)
	id := rand.Intn(RAND_MAX)
	ps, prs := psTopic[topic]
	if prs == false {
		details := map[string]string{}
		uri := "wamp.error.topic_not_exist"
		resp = []interface{}{Error, Publish, request, details, uri}
	} else {
		details := map[string]string{}
		event := []interface{}{Event, ps.id, id, details}
		if len(msg) > 4 {
			args := msg[4]
			event = append(event, args)
		}
		if len(msg) > 5 {
			kwargs := msg[5]
			event = append(event, kwargs)
		}
		log.Debugln("Publish from", se.id, "event", ps.topic)
		for _, subscriber := range ps.subscribers {
			if subscriber != se {
				go subscriber.Conn.Write(event)
			}
		}
		resp = []interface{}{Published, request, id}
	}
	return
}
