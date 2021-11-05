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

func (f *cinFlow) findArg(label string) (retChoice string, retMsg json.RawMessage) {
	if i, cnt := f.bestIndex, len(f.params); i < cnt && f.params[i].label == label {
		retChoice, retMsg = f.getArg(i)
	} else if i == 0 && len(label) == 0 {
		// in the case of anonymous swaps -- the label we parsed is actually a choice.
		// ex. "EventBlock kind:"
		retChoice, retMsg = f.params[i].label, f.args[i]
	} else {
		for ; i < cnt; i++ {
			if n := f.params[i]; n.label == label {
				retChoice, retMsg = f.getArg(i)
				break
			}
		}
	}
	return
}

func (f *cinFlow) getArg(i int) (string, json.RawMessage) {
	f.bestIndex = i + 1 // next time we're most likely on the next arg.
	return f.params[i].choice, f.args[i]
}
