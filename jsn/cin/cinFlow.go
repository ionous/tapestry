package cin

import (
	r "reflect"

	"github.com/ionous/errutil"
)

type cinFlow struct {
	name      string
	params    []Parameter
	args      r.Value // a slice of interfaces
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

// returns the argument for the passed parameter label;
// nothing if not found; errors if known to be invalid.
func (f *cinFlow) getArg(key string) (ret r.Value, err error) {
	if i := f.getParamIndex(key); i >= 0 {
		if c := f.params[i].Choice; len(c) > 0 {
			err = errutil.Fmt("expected no choice for key %q have %q", key, c)
		} else {
			ret = f.args.Index(i).Elem()
			f.bestIndex = i + 1
		}
	}
	return
}

func (f *cinFlow) getPick(key string) (retMsg r.Value, retChoice string, err error) {
	// the signature parser can't distinguish b/t a leading first selector, and a leading anonymous choice:
	//  ex. "Command choice:" -- so the param array has the choice in the label's spot
	if len(key) == 0 && f.bestIndex == 0 && len(f.params) > 0 && len(f.params[0].Choice) == 0 {
		retChoice = f.params[0].Label
		retMsg = f.args.Index(0).Elem()
		f.bestIndex = 1
	} else {
		// otherwise we expect named selector/choice pairs
		// "Command selector choice:"
		if i := f.getParamIndex(key); i >= 0 {
			if c := f.params[i].Choice; len(c) == 0 {
				err = errutil.Fmt("expected a trailing choice for key %q", key)
			} else {
				retChoice = c
				retMsg = f.args.Index(i).Elem()
				f.bestIndex = i + 1
			}
		}
	}
	return
}

// returns the index of the labeled parameter, or -1 if not found.
// only looks at or to the right of the current "bestIndex"
func (f *cinFlow) getParamIndex(label string) (ret int) {
	ret = -1 // provisionally
	for i, cnt := f.bestIndex, len(f.params); i < cnt; i++ {
		if n := f.params[i]; n.Label == label {
			ret = i
			break
		}
	}
	return
}
