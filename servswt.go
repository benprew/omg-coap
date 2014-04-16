package coap

import (
	"math/rand"
	"time"
)

// Observer of some resource
type observer struct {
	tid     uint16
	rsc     string
	addr    *RemoteAddr
	send    chan *transmission
	ack     chan uint16
	payload string
}

type transmission struct {
	ctx string
	id  string
	msg *Message
	ret *chan bool
}

// Internal switch
type coapswitch struct {
	// switch operations
	register   chan *observer
	unregister chan *observer
	// networking operations
	ack       chan map[string]uint16
	unicast   chan *transmission
	broadcast chan *transmission
	// observers cache by context/id
	observers map[string]map[string]*observer
}

var s = coapswitch{
	ack:        make(chan map[string]uint16),
	register:   make(chan *observer),
	unregister: make(chan *observer),

	unicast:   make(chan *transmission),
	broadcast: make(chan *transmission),

	observers: make(map[string]map[string]*observer),
}

// Observer transmitter.
func (o *observer) transmitter() {
	/*
			   Implementation Note:  Several implementation strategies can be
			        employed for generating Message IDs.  In the simplest case a CoAP
			        endpoint generates Message IDs by keeping a single Message ID
			        variable, which is changed each time a new Confirmable or Non-
			        confirmable message is sent regardless of the destination address
			        or port.  Endpoints dealing with large numbers of transactions
			        could keep multiple Message ID variables, for example per prefix
			        or destination address (note that some receiving endpoints may not
			        be able to distinguish unicast and multicast packets addressed to
			        it, so endpoints generating Message IDs need to make sure these do
			        not overlap).  It is strongly recommended that the initial value
			        of the variable (e.g., on startup) be randomized, in order to make
			        successful off-path attacks on the protocol less likely.
		     http://tools.ietf.org/html/draft-ietf-core-coap-18#section-4.4
	*/
	o.tid = uint16(rand.Intn(65536))
	for t := range o.send {
		o.tid++
		go func(m Message, tid uint16, maxRetransmit int, responseTimeout time.Duration, alive chan bool) {
			m.MessageID = tid
			ticker := time.NewTicker(responseTimeout)
			for i := 0; i <= maxRetransmit; i++ {
				if i == 0 {
					debugMsg("** transmission of message %v of [%s] resource to %s", m.MessageID, o.rsc, o.addr)
					Transmit(o.addr, m)
				}
				select {
				case a := <-o.ack:
					if a == m.MessageID {
						alive <- true
						return
					}
				case <-ticker.C:
					if i == maxRetransmit {
						debugMsg("** transmission of message %v of [%s] resource timeout", m.MessageID, o.rsc)
						alive <- false
						s.unregister <- o
						return
					}
					if i <= maxRetransmit {
						debugMsg("** retransmission #%d of message %v of [%s] resource to %s", i+1, m.MessageID, o.rsc, o.addr)
						Transmit(o.addr, m)
					}
				}
			}
		}(*t.msg, o.tid, MaxRetransmit, ResponseTimeout, *t.ret)
	}
}

// Observe switch runtime.
func (s *coapswitch) run() {
	for {
		select {
		case o := <-s.register:
			obs, ok := s.observers[o.rsc]
			if !ok {
				obs = make(map[string]*observer)
				s.observers[o.rsc] = obs
			}
			s.observers[o.rsc][o.payload] = o
			debugMsg("** observer %s added in [%s]", o.addr, o.rsc)
		case o := <-s.unregister:
			delete(s.observers[o.rsc], o.payload)
			close(o.ack)
			close(o.send)
			debugMsg("** observer %s of [%s] removed", o.addr, o.rsc)
		case a := <-s.ack:
			// ack have transmission id only, match by addr
			for _, obs := range s.observers {
				for _, o := range obs {
					// send ack to right observer
					if tid, ok := a[o.addr.String()]; ok && tid != 0 {
						o.ack <- tid
					}
				}
			}
		case t := <-s.broadcast:
			// message to all of context
			obs, alive := s.observers[t.ctx]
			if alive {
				for _, o := range obs {
					o.send <- t
				}
			}
		case t := <-s.unicast:
			// message to observer of context with specific id
			o, alive := s.observers[t.ctx][t.id]
			if alive {
				o.send <- t
			} else {
				close(*t.ret)
			}
		}
	}
}
