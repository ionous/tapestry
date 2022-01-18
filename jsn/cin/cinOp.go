package cin

import (
	"encoding/json"
	"strings"

	"github.com/ionous/errutil"
)

type Op struct {
	Key string          // unparsed signature
	Cmt string          // comment from "--" fields
	Msg json.RawMessage // probably an array of values
}

// ex. {"Story:":  [...]}
// except note that literals are stored as a single literal value;
// functions with one parameter dont use the array, and
// functions without parameters are stored as simple strings.
func ReadOp(msg json.RawMessage) (ret Op, err error) {
	var d map[string]json.RawMessage
	if e := json.Unmarshal(msg, &d); e == nil {
		ret, err = parseOp(d) // start by trying to read the {} format
	} else {
		err = json.Unmarshal(msg, &ret.Key) // then try a raw string for parameterless commands.
	}
	return
}

func (op *Op) ReadMsg() (retSig Signature, retArgs []json.RawMessage, err error) {
	// we allow ( require really ) single arguments to be stored directly
	// rather than embedded in an array
	// to make it optional, we'd really need a parallel parser to attempt to interpret the argument bytes in multiple ways.
	if sig, e := ReadSignature(op.Key); e != nil {
		err = e
	} else {
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
		} else if len(out.Key) > 0 {
			err = errutil.New("expected only a single key", d)
			break
		} else {
			out.Key, out.Msg = k, v
			continue // keep going to catch errors
		}
	}
	if err == nil {
		// in the case that there was no command but there was a comment marker
		// let the command *be* the comment marker
		if len(out.Key) == 0 && hadComment {
			out.Key = commentMarker
		}
		ret = out
	}
	return
}

const commentMarker = "--"
