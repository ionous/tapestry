package cin

import (
	r "reflect"
	"strings"

	"git.sr.ht/~ionous/tapestry/rt/markup"
	"github.com/ionous/errutil"
)

// Op represents a partial decoding of a command in the compact format.
// ex. {"Sig:":  <Msg>, "--": "Markup"}
// ( note: the registry uses the raw signature without additional processing to find associated golang struct )
type Op struct {
	Sig    string         // raw signature containing one or more colon separators
	Markup map[string]any // metadata from "--" fields
	Msg    r.Value        // probably an array of values
}

// ReadOp - interpret the passed object as the start of a compact command.
func ReadOp(msg r.Value) (ret Op, err error) {
	if t := msg.Type(); !IsValidMap(t) {
		err = errutil.Fmt("expected a compact command, not %s", t)
	} else {
		ret, err = parseOp(msg)
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
// Errors if the number of separators in its sig differs from the number of parameters in its msg.
// retArgs is guaranteed to be a slice
func (op *Op) ReadMsg() (retSig Signature, retArgs r.Value, err error) {
	if sig, e := ReadSignature(op.Sig); e != nil {
		err = e
	} else {
		// we *require* that single arguments get stored directly ( rather than embedded in an array. )
		// ideally it would be optional, but it gets a bit weird when the single argument is itself an array.
		// there's no way to distinguish that case without knowing the desired format of the argument ( and we dont here. )
		// fix? maybe the cout should always use an array.... it just seemed verbose at the time.
		var args r.Value
		switch pn := len(sig.Params); pn {
		case 0:
			// tbd: this used to return a nil slice;
			// now it returns invalid. is that okay?
		case 1:
			args = r.ValueOf([]any{op.Msg.Interface()})
		default:
			slice := op.Msg
			if t := slice.Type(); !IsValidSlice(t) {
				err = errutil.Fmt("expected a slice of arguments, not a(n) %s", t)
			} else if an := slice.Len(); an != pn {
				err = errutil.Fmt("expected %s with %d args, has %d args",
					sig.DebugString(), pn, an)
			}
			args = op.Msg
		}
		if err == nil {
			retSig = sig
			retArgs = args
		}
	}
	return
}

// expects a map of string to value
func parseOp(obj r.Value) (ret Op, err error) {
	var out Op
	for it := obj.MapRange(); it.Next(); {
		k, v := it.Key().String(), it.Value().Elem()
		if strings.HasPrefix(k, markupMarker) {
			if key := k[len(markupMarker):]; len(key) == 0 {
				out.AddMarkup(markup.Comment, v.Interface())
			} else {
				out.AddMarkup(key, v.Interface())
			}
		} else if len(out.Sig) > 0 {
			err = errutil.New("expected only a single key")
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
			if _, ok := out.Markup[markup.Comment]; ok {
				out.Sig = markupMarker
			}
		}
		ret = out
	}
	return
}

const markupMarker = "--"
