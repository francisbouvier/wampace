package wamp

import log "github.com/Sirupsen/logrus"

type sub struct {
	id           int
	subRequest   int
	unsubRequest int
	topic        string
	Subscribed   chan error
	Unsubscribed chan error
	cbk          callback
}

type pub struct {
	id          int
	request     int
	topic       string
	subscribers map[int]*Session
	Published   chan error
	a           args
	k           kwargs
}

func (ps *pub) addSub(se *Session) {
	if _, prs := ps.subscribers[se.id]; prs == false {
		ps.subscribers[se.id] = se
		log.Debugln("Add subscriber", se.id, "to topic", ps.topic)
	}
	return
}

func (ps *pub) removeSub(se *Session) {
	if _, prs := ps.subscribers[se.id]; prs == true {
		delete(ps.subscribers, se.id)
		log.Debugln("Remove subscriber", se.id, "from topic", ps.topic)
	}
	return
}
