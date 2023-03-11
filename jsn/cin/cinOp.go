package cin

import (
	"encoding/json"
	"strings"

	"github.com/ionous/errutil"
)

// Op represents a partial decoding of a command in the compact format.
// ex. {"Sig:":  <Msg>, "--": "Markup"}
// ( note: the registry uses the raw signature without additional processing to find associated golang struct )
type Op struct {
	Sig    string          // raw signature containing one or more colon separators
	Markup map[string]any  // metadata from "--" fields
	Msg    json.RawMessage // probably an array of values
}

// ReadOp - interpret the passed json as the start of a compact command.
// fix? replace this with a custom implementation of json.Unmarshaler ( UnmarshalJSON() )
func ReadOp(msg json.RawMessage) (ret Op, err error) {
	var d map[string]json.RawMessage
	if e := json.Unmarshal(msg, &d); e != nil {
		err = e
	} else {
		ret, err = parseOp(d) // start by trying to read the {} format
	}
	return
}

func (op *Op) AddMarkup(k string, value any) {
	if op.Markup == nil {
		op.Markup = make(map[string]any)
	}
	op.Markup[k] = value
}

// ReadMsg - given a valid Op, split out its call signature and associated parameters.
// Errors if the number of separators in its sig differs from the the number of parameters in its msg.
func (op *Op) ReadMsg() (retSig Signature, retArgs []json.RawMessage, err error) {
	if sig, e := ReadSignature(op.Sig); e != nil {
		err = e
	} else {
		// we *require* that single arguments get stored directly ( rather than embedded in an array. )
		// ideally it would be optional, but it gets a bit weird when the single argument is itself an array.
		// there's no way to distinguish that case without knowing the desired format of the argument ( and we dont here. )
		// fix? maybe the cout should always use an array.... it just seemed verbose at the time.
		var args []json.RawMessage
		pn := len(sig.Params)
		if pn == 1 {
			args = []json.RawMessage{op.Msg}
		} else if pn > 1 {
			err = json.Unmarshal(op.Msg, &args)
		}
		if err == nil {
			if an := len(args); pn != an {
				snippet := "???"
				if an > 0 {
					x := args[0]
					if cap := 25; len(x) > cap {
						x = x[:cap]
					}
					snippet = string(x)
				}
				err = errutil.Fmt("%q given %d args: %s", sig.DebugString(), an, snippet)
			} else {
				retSig = sig
				retArgs = args
			}
		}
	}
	return
}

func parseOp(d map[string]json.RawMessage) (ret Op, err error) {
	var out Op
	for k, v := range d {
		if strings.HasPrefix(k, markupMarker) {
			var value any
			if e := json.Unmarshal(v, &value); e != nil {
				err = errutil.New("couldnt read markup at", k, e)
			} else if key := k[len(markupMarker):]; len(key) == 0 {
				out.AddMarkup("comment", value)
			} else {
				out.AddMarkup(key, value)
			}
		} else if len(out.Sig) > 0 {
			err = errutil.New("expected only a single key", d)
			break
		} else {
			out.Sig, out.Msg = k, v
			continue // keep going to catch errors
		}
	}
	if err == nil {
		// in the case that there was no command but there was a comment marker
		// let the command *be* the comment marker
		if len(out.Sig) == 0 {
			if _, ok := out.Markup["comment"]; ok {
				out.Sig = markupMarker
			}
		}
		ret = out
	}
	return
}

const markupMarker = "--"
