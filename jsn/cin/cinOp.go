package cin

import (
	"strings"

	"git.sr.ht/~ionous/tapestry/lang/markup"
	"github.com/ionous/errutil"
)

// Op represents a partial decoding of a command in the compact format.
// ex. {"Sig:":  <Msg>, "--": "Markup"}
// ( note: the registry uses the raw signature without additional processing to find associated golang struct )
type Op struct {
	Sig    string         // raw signature containing one or more colon separators
	Markup map[string]any // metadata from "--" fields
	Msg    any            // usually an array of values, or a single value shortcut for an array.
}

// ReadOp - interpret the passed object as the start of a compact command.
func ReadOp(msg map[string]any) (ret Op, err error) {
	return parseOp(msg)
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
func (op *Op) ReadMsg() (retSig Signature, retArgs []any, err error) {
	if sig, e := ReadSignature(op.Sig); e != nil {
		err = e
	} else {
		// we *require* that single arguments get stored directly ( rather than embedded in an array. )
		// ideally it would be optional, but it gets a bit weird when the single argument is itself an array.
		// there's no way to distinguish that case without knowing the desired format of the argument ( and we dont here. )
		// fix? maybe the cout should always use an array.... it just seemed verbose at the time.
		var args []any
		switch pn, body := len(sig.Params), op.Msg; pn {
		case 0:
			// tbd: this used to return a nil slice;
			// now it returns invalid. is that okay?
		case 1:
			args = []any{body}
		default:
			if slice, ok := body.([]any); !ok {
				err = errutil.Fmt("expected a slice of arguments, not a(n) %T", body)
			} else if an := len(slice); an != pn {
				err = errutil.Fmt("expected %s with %d args, has %d args",
					sig.DebugString(), pn, an)
			} else {
				args = slice
			}
		}
		if err == nil {
			retSig = sig
			retArgs = args
		}
	}
	return
}

// expects a map of string to value
func parseOp(msg map[string]any) (ret Op, err error) {
	var out Op
	for k, v := range msg {
		if strings.HasPrefix(k, markupMarker) {
			if key := k[len(markupMarker):]; len(key) == 0 {
				out.AddMarkup(markup.Comment, v)
			} else {
				out.AddMarkup(key, v)
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
		// fix: remove. standalone comments should be the comment command only
		// ( jsnTest still has some of these. )
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
