package coap

import (
	"math/rand"
	"time"
)

// Observer of some resource
type Observer struct {
	tid  uint16
	rsc  string
	Addr *RemoteAddr
	Send chan *Message
	ack  chan uint16
  payload string
}

// Internal hub
type HUB struct {
	ack        chan map[string]uint16
	register   chan *Observer
	unregister chan *Observer
	Observers  map[string]map[string]*Observer
}

var H = HUB{
	ack:        make(chan map[string]uint16),
	register:   make(chan *Observer),
	unregister: make(chan *Observer),
	Observers:  make(map[string]map[string]*Observer),
}

// Observer transmitter.
func (o *Observer) transmitter() {
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
	for msg := range o.Send {
		o.tid++
		go func(m Message, tid uint16, mr int, rt time.Duration) {
			m.MessageID = tid
			ticker := time.NewTicker(rt)
			for i := 0; i <= mr; i++ {
				if i == 0 {
					debugMsg("** transmission of message %v of [%s] resource to %s", m.MessageID, o.rsc, o.Addr)
					Transmit(o.Addr, m)
				}
				select {
				case a := <-o.ack:
					if a == m.MessageID {
						return
					}
				case <-ticker.C:
					if i == mr {
						debugMsg("** transmission of message %v of [%s] resource timeout", m.MessageID, o.rsc)
						H.unregister <- o
						return
					}
					if i <= mr {
						debugMsg("** retransmission #%d of message %v of [%s] resource to %s", i+1, m.MessageID, o.rsc, o.Addr)
						Transmit(o.Addr, m)
					}
				}
			}
		}(*msg, o.tid, MaxRetransmit, ResponseTimeout)
	}
}

// Observe hub runtime.
func (H *HUB) run() {
	for {
		select {
		case o := <-H.register:
			obs, ok := H.Observers[o.rsc]
			if !ok {
				obs = make(map[string]*Observer)
				H.Observers[o.rsc] = obs
			}
			H.Observers[o.rsc][o.payload] = o
			debugMsg("** observer %s added in [%s]", o.Addr, o.rsc)
		case o := <-H.unregister:
			delete(H.Observers[o.rsc], o.payload)
			close(o.ack)
			close(o.Send)
			debugMsg("** observer %s of [%s] removed", o.Addr, o.rsc)
		case a := <-H.ack:
			for _, o := range H.Observers {
				for _, obs := range o {
					// send ack to right observer
					if tid, ok := a[obs.Addr.String()]; ok && tid != 0 {
						obs.ack <- tid
					}
				}
			}
		}
	}
}
