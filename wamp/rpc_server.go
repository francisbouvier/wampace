package wamp

import (
	"math/rand"
)

func (s Server) registred(se *Session, msg []interface{}) (resp []interface{}, err error) {
	request := int(msg[1].(float64))
	procedure := msg[3].(string)
	var rp rProc
	rpcLock.Lock()
	_, prs := rps[procedure]
	if prs == false {
		rp = rProc{
			id:        rand.Intn(RAND_MAX),
			procedure: procedure,
		}
		rps[procedure] = &rp
		rpsID[rp.id] = &rp
		rp.addCallee(se)
		resp = []interface{}{Registred, request, rp.id}
	} else {
		details := map[string]string{}
		uri := "wamp.error.procedure_already_exists"
		resp = []interface{}{Error, Register, request, details, uri}
	}
	rpcLock.Unlock()
	return
}

func (s Server) unregistred(se *Session, msg []interface{}) (resp []interface{}, err error) {
	request := int(msg[1].(float64))
	id := int(msg[2].(float64))
	rp, prs := rpsID[id]
	if prs == true {
		rp.removeCallee(se)
		resp = []interface{}{Unregistred, request}
	} else {
		details := map[string]string{}
		uri := "wamp.error.no_such_registration"
		resp = []interface{}{Error, Unregister, request, details, uri}
	}
	return
}

func (s Server) invocation(se *Session, msg []interface{}) (resp []interface{}, err error) {
	request := int(msg[1].(float64))
	procedure := msg[3].(string)
	rp, prs := rps[procedure]
	if prs == false {
		details := map[string]string{}
		uri := "wamp.error.no_such_procedure"
		resp = []interface{}{Error, Call, request, details, uri}
	} else {
		rc := RCall{
			id:      rand.Intn(RAND_MAX),
			request: request,
			rp:      rp,
			caller:  se,
		}
		rcsID[rc.id] = &rc
		details := map[string]string{}
		call := []interface{}{Invocation, rc.id, rp.id, details}
		if len(msg) > 4 {
			args := msg[4]
			call = append(call, args)
		}
		if len(msg) > 5 {
			kwargs := msg[5]
			call = append(call, kwargs)
		}
		go rp.callee.Conn.Write(call)
	}
	return
}

func (s Server) result(w *Session, msg []interface{}) (resp []interface{}, err error) {
	id := int(msg[1].(float64))
	rc, prs := rcsID[id]
	if prs == true {
		details := map[string]string{}
		yield := []interface{}{Result, rc.request, details}
		if len(msg) > 3 {
			args := msg[3]
			yield = append(yield, args)
		}
		if len(msg) > 4 {
			kwargs := msg[4]
			yield = append(yield, kwargs)
		}
		go rc.caller.Conn.Write(yield)
	}
	return
}

func (s Server) callError(se *Session, msg []interface{}) (resp []interface{}, err error) {
	id := int(msg[2].(float64))
	rc, prs := rcsID[id]
	if prs == true {
		details := msg[3]
		uri := msg[4].(string)
		yield := []interface{}{Error, Call, rc.request, details, uri}
		if len(msg) > 5 {
			args := msg[5]
			yield = append(yield, args)
		}
		if len(msg) > 6 {
			kwargs := msg[6]
			yield = append(yield, kwargs)
		}
		go rc.caller.Conn.Write(yield)
	}
	return
}
