package assign

import (
	"fmt"

	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"git.sr.ht/~ionous/tapestry/support/inflect"
)

// turn a series of assignments ( FromX commands ) into a slice of arguments.
func MakeArgs(as ...rt.Assignment) (ret []Arg) {
	for _, a := range as {
		ret = append(ret, Arg{Value: a})
	}
	return
}
func ExpandArgs(run rt.Runtime, args []Arg) (retKeys []string, retVals []rt.Value, err error) {
	if len(args) > 0 {
		keys, vals := make([]string, 0, len(args)), make([]rt.Value, len(args))
		for i, a := range args {
			if val, e := safe.GetAssignment(run, a.Value); e != nil {
				err = fmt.Errorf("%w while reading arg %d(%s)", e, i, a.Name)
				break
			} else if n := inflect.Normalize(a.Name); len(n) > 0 {
				keys = append(keys, n)
				vals[i] = val
			} else if len(keys) > 0 {
				err = fmt.Errorf("unnamed arguments must precede all named arguments %d", i)
			} else {
				vals[i] = val
			}
		}
		if err == nil {
			retKeys, retVals = keys, vals
		}
	}
	return
}
