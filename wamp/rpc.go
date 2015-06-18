package wamp

import log "github.com/Sirupsen/logrus"

type RCall struct {
	id      int
	request int
	Result  chan error
	Args    args
	Kwargs  kwargs
	rp      *rProc
	caller  *Session
}

type rProc struct {
	id           int
	regRequest   int
	unregRequest int
	procedure    string
	Registred    chan error
	Unregistred  chan error
	cbk          proc
	callee       *Session
}

func (rp *rProc) addCallee(se *Session) {
	rp.callee = se
	log.Debugln("Register callee", se.id, "with procedure", rp.procedure)
	return
}

func (rp *rProc) removeCallee(se *Session) {
	if se == rp.callee {
		delete(rps, rp.procedure)
		delete(rpsID, rp.id)
		log.Debugln("Unregister callee", se.id, "with procedure", rp.procedure)
	}
	return
}
