package frame

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/dl/debug"
	"git.sr.ht/~ionous/tapestry/lang/encode"
	"git.sr.ht/~ionous/tapestry/qna/query"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/print"
)

func NewShuttle(run rt.Runtime, dec *query.QueryDecoder) *Shuttle {
	c := &Shuttle{
		run: run,
		dec: dec,
	}
	note := rt.Notifier{
		StartedScene:    c.out.onStartScene,
		EndedScene:      c.out.onEndScene,
		ChangedState:    c.out.onChangeState,
		ChangedRelative: c.out.onChangeRel,
	}
	run.SetNotifier(note)
	run.SetWriter(print.NewLineSentences(&c.out.buf))
	return c
}

// Shuttle uses json commands to talk back and forth to a runtime.
// see: idl frame.
type Shuttle struct {
	dec *query.QueryDecoder // used to decode queries
	run rt.Runtime
	out Collector
}

func (c *Shuttle) Restart(scene string) (err error) {
	// FIX: tear down.
	run := c.run
	run.ActivateDomain("")
	// run.Survey().GetFocalObject()}
	focus := rt.StringOf("self")
	//
	if e := run.ActivateDomain(scene); e != nil {
		err = e
	} else {
		_, err = run.Call("start game", affine.None, nil, []rt.Value{focus})
	}
	return
}

func (c *Shuttle) Post(w io.Writer, endpoint string, msg json.RawMessage) (err error) {
	switch endpoint {

	case Restart.String():
		var scene string
		if e := json.Unmarshal(msg, &scene); e != nil {
			err = errors.New("invalid scene")
		} else if e := c.Restart(scene); e != nil {
			err = e
		} else {
			evts := c.out.GetEvents()
			err = writeFrames(w, []Frame{{Events: evts}})
		}

	case Query.String():
		var qs []json.RawMessage
		if e := json.Unmarshal(msg, &qs); e != nil {
			err = errors.New("invalid query")
		} else if frames, e := c.readFrames(qs); e != nil {
			err = e
		} else {
			err = writeFrames(w, frames)
		}

	default:
		err = fmt.Errorf("unknown endpoint %s", endpoint)
	}
	return
}

func (c *Shuttle) readFrames(qs []json.RawMessage) (ret []Frame, err error) {
	for _, q := range qs {
		var f Frame // a frame contains many events.
		if a, e := c.dec.DecodeAssignment(affine.None, q); e != nil {
			err = e
			break
		} else if v, e := a.GetAssignedValue(c.run); e == nil && v != nil {
			f.Result = debug.Stringify(v)
			// success!
		} else if e != nil {
			//
			var sig rt.Signal
			if errors.As(e, &sig) {
				// fix? this is a little wonky
				// signals should probably be a first class method in the runtime?
				c.out.onGameEvent(sig)
			} else {
				f.Error = e.Error()
				// return this error as part of the frame
			}
		}
		//
		f.Events = c.out.GetEvents()
		ret = append(ret, f)
	}
	return
}

// write the list of frames to the output as raw json.
func writeFrames(w io.Writer, frames Frame_Slice) (err error) {
	var enc encode.Encoder
	if d, e := enc.Encode(&frames); e != nil {
		err = e
	} else {
		js := json.NewEncoder(w)
		js.SetEscapeHTML(false)
		err = js.Encode(d)
	}
	return
}
