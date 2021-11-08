package cin

import (
	"encoding/json"

	"github.com/ionous/errutil"
)

type cinFlow struct {
	params    []sigParam
	args      []json.RawMessage
	bestIndex int
}

func newFlowData(k string, msg json.RawMessage) (ret *cinFlow, err error) {
	var sig sigReader
	var args []json.RawMessage
	sig.readSig(k)
	// we allow ( require really ) single arguments to be stored directly
	// rather than embedded in an array
	// to make it optional, we'd really need a parallel parser to attempt to interpret the argument bytes in multiple ways.
	pn := len(sig.params)
	if pn == 1 {
		args = []json.RawMessage{msg}
	} else if pn > 1 {
		err = json.Unmarshal(msg, &args)
	}
	if err == nil {
		if an := len(args); pn != an {
			err = errutil.New("mismatched params and args", pn, an)
		} else {
			ret = &cinFlow{params: sig.params, args: args}
		}
	}
	return
}

func (f *cinFlow) getArg(key string) (ret json.RawMessage, err error) {
	if i := f.getArgAt(key); i >= 0 {
		if c := f.params[i].choice; len(c) > 0 {
			err = errutil.Fmt("expected no choice for key %q have %q", key, c)
		} else {
			ret = f.args[i]
			f.bestIndex = i + 1
		}
	}
	return
}
func (f *cinFlow) getPick(key string) (retMsg json.RawMessage, retChoice string, err error) {
	// the signature parser can't distinguish b/t a leading first selector, and a leading anonymous choice:
	//  ex. "Command choice:" -- so the param array has the choice in the label's spot
	if len(key) == 0 && f.bestIndex == 0 && len(f.params) > 0 && len(f.params[0].choice) == 0 {
		retChoice = f.params[0].label
		retMsg = f.args[0]
		f.bestIndex = 1
	} else {
		// otherwise we expect named selector/choice pairs
		// "Command selector choice:"
		if i := f.getArgAt(key); i >= 0 {
			if c := f.params[i].choice; len(c) == 0 {
				err = errutil.Fmt("expected a trailing choice for key %q", key)
			} else {
				retChoice = c
				retMsg = f.args[i]
				f.bestIndex = i + 1
			}
		}
	}
	return
}

func (f *cinFlow) getArgAt(key string) (ret int) {
	ret = -1 // provisionally
	for i, cnt := f.bestIndex, len(f.params); i < cnt; i++ {
		if n := f.params[i]; n.label == key {
			ret = i
			break
		}
	}
	return
}
