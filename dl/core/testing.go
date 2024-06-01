package core

import (
	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/dl/logic"
	"git.sr.ht/~ionous/tapestry/rt"
)

// turn a series of assignments ( FromX commands ) into a slice of arguments.
func MakeArgs(as ...rt.Assignment) (ret []assign.Arg) {
	for _, a := range as {
		ret = append(ret, assign.Arg{Value: a})
	}
	return
}

// for tests
func MakeRule(filter rt.BoolEval, exe ...rt.Execute) (ret rt.Rule) {
	if filter == nil {
		ret.Exe = exe
	} else {
		ret = rt.Rule{Exe: []rt.Execute{
			&logic.ChooseBranch{
				Condition: filter,
				Exe:       exe,
			}}}
	}
	return
}
