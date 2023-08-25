package cin

import (
	"encoding/json"

	"github.com/ionous/errutil"
)

type cinFlow struct {
	name      string
	params    []Parameter
	args      []json.RawMessage
	bestIndex int
}

func newFlowData(op Op) (ret *cinFlow, err error) {
	if sig, args, e := op.ReadMsg(); e != nil {
		err = e
	} else {
		ret = &cinFlow{name: sig.Name, params: sig.Params, args: args}
	}
	return
}

func (f *cinFlow) getArg(key string) (ret json.RawMessage, err error) {
	if i := f.getArgAt(key); i >= 0 {
		if c := f.params[i].Choice; len(c) > 0 {
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
	if len(key) == 0 && f.bestIndex == 0 && len(f.params) > 0 && len(f.params[0].Choice) == 0 {
		retChoice = f.params[0].Label
		retMsg = f.args[0]
		f.bestIndex = 1
	} else {
		// otherwise we expect named selector/choice pairs
		// "Command selector choice:"
		if i := f.getArgAt(key); i >= 0 {
			if c := f.params[i].Choice; len(c) == 0 {
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
		if n := f.params[i]; n.Label == key {
			ret = i
			break
		}
	}
	return
}
