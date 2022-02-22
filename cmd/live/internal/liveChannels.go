package live

import (
	"io"
	"log"
	"sync"

	"git.sr.ht/~ionous/tapestry/dl/play"
	"git.sr.ht/~ionous/tapestry/web/js"
)

// see NewChannels
type Channels struct {
	msgs  chan play.PlayMessage
	input chan string
	// only one reader at a time to ensure messages dont wind up interleaved in multiple response.
	// though we could support storing all the messages ( in memory or on disk )
	// along with "windowing" requests from the client if it were really useful at all.
	msgReader sync.Mutex
}

func NewChannels() *Channels {
	return &Channels{
		msgs:  make(chan play.PlayMessage, 1024),
		input: make(chan string, 8),
	}
}

func (cs *Channels) ChangeMode(n string) {
	cs.msgs <- &play.PlayMode{Mode: play.PlayModes{Str: n}}
}

func (cs *Channels) Fatal(e error) {
	cs.ChangeMode(play.PlayModes_Error)
	cs.msgs <- &play.PlayLog{Log: e.Error()}
}

func (cs *Channels) WriteMessages(w io.Writer, wait bool) (err error) {
	// read all the messages into a json block.
	cs.msgReader.Lock()
	cnt, str := readMsgs(cs.msgs, wait)
	cs.msgReader.Unlock()
	//
	if _, e := io.WriteString(w, str); e != nil {
		err = e
	} else {
		log.Println("wrote", cnt, "msgs")
	}
	return
}

// generate json from any pending messages in the passed channel
func readMsgs(msgs <-chan play.PlayMessage, wait bool) (retCnt int, retStr string) {
	var out js.Builder // fix: really need to make this streamable.

	out.Brace(js.Array, func(a *js.Builder) {
		// some arbitrary maximum size
		for i := 0; i < 512; i++ {
			if msg, ok := readMsg(msgs, wait); !ok {
				break
			} else {
				if i > 0 {
					a.R(js.Comma)
				}
				str, _ := Marshal(msg)
				a.S(str)
				retCnt++
				wait = false
			}
		}
	})
	retStr = out.String()
	return
}

func readMsg(msgs <-chan play.PlayMessage, wait bool) (ret play.PlayMessage, okay bool) {
	if wait {
		ret = <-msgs
		okay = true
	} else {
		select {
		case msg := <-msgs:
			ret = msg
			okay = true
		default:
		}
	}
	return
}
