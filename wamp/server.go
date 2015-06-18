package wamp

import (
	"sync"
)

var (
	rps     = map[string]*rProc{}
	rpsID   = map[int]*rProc{}
	rcsID   = map[int]*RCall{}
	rpcLock sync.Mutex
)

type Server struct {
}

func (s Server) welcome(se *Session, msg []interface{}) (resp []interface{}, err error) {
	realm := msg[1].(string)
	se.open(realm)
	details := map[string]map[string]map[string]interface{}{
		"roles": {
			"broker": {},
			"dealer": {},
		},
	}
	resp = []interface{}{Welcome, se.id, details}
	return
}

func (s Server) abort(se *Session, msg string, uri string) (resp []interface{}, err error) {
	se.Close()
	details := map[string]string{
		"message": msg,
	}
	resp = []interface{}{Abort, details, uri}
	return
}

func (s Server) goodbye(se *Session, msg []interface{}) (resp []interface{}, err error) {
	se.Close()
	details := map[string]string{}
	reason := "wamp.error.goodbye_and_out"
	resp = []interface{}{Goodbye, details, reason}
	return
}

func (s Server) Route(se *Session, msg []interface{}) (resp []interface{}, err error) {
	code := int(msg[0].(float64))
	if code == Hello {
		// Trying to connect a peer already connected cause abort
		if err = se.check(); err == nil {
			uri := "wamp.error.peer_already_connected"
			return s.abort(se, "Peer already connected", uri)
		}
		return s.welcome(se, msg)
	} else {
		// Sending message from a peer not connected cause abort (general error ?)
		if err = se.check(); err != nil {
			return s.abort(se, err.Error(), "wamp.error.peer_not_connected")
		}
		switch code {
		case Goodbye:
			return s.goodbye(se, msg)
		case Subscribe:
			return s.subscribed(se, msg)
		case Unsubscribe:
			return s.unsubscribed(se, msg)
		case Publish:
			return s.published(se, msg)
		case Register:
			return s.registred(se, msg)
		case Unregister:
			return s.unregistred(se, msg)
		case Call:
			return s.invocation(se, msg)
		case Yield:
			return s.result(se, msg)
		case Error:
			sub := int(msg[0].(float64))
			switch sub {
			case Invocation:
				return s.callError(se, msg)
			}
		default:
			return s.abort(se, "Unknow message", "wamp.error.unknow_message")
		}
	}
	return
}

func NewServer() (s Server) {
	s = Server{}
	return
}
