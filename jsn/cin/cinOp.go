package cin

import (
	"encoding/json"
	"strings"

	"github.com/ionous/errutil"
)

// Op represents a partial decoding of a command in the compact format.
// ex. {"Sig:":  <Msg>, "--": "Cmt"}
// ( note: the registry uses the raw signature without additional processing to find associated golang struct )
type Op struct {
	Sig string          // raw signature containing one or more colon separators
	Cmt string          // comment from "--" fields
	Msg json.RawMessage // probably an array of values
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
				err = errutil.New("mismatched params and args", pn, an)
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
	var hadComment bool
	for k, v := range d {
		if k == commentMarker {
			hadComment = true
			if e := json.Unmarshal(v, &out.Cmt); e != nil {
				var lines []string
				if e := json.Unmarshal(v, &lines); e != nil {
					err = errutil.New("couldnt read comment", e)
					break
				}
				out.Cmt = strings.Join(lines, "\n")
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
		if len(out.Sig) == 0 && hadComment {
			out.Sig = commentMarker
		}
		ret = out
	}
	return
}

const commentMarker = "--"
