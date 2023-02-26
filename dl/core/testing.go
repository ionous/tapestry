package core

import (
	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/rt"
)

// turn a series of statements into a slice of execute statements.
func MakeActivity(exe ...rt.Execute) []rt.Execute {
	return exe
}

// turn a series of assignments ( FromX commands ) into a slice of arguments.
func MakeArgs(as ...assign.Assignment) (ret []assign.Arg) {
	for _, a := range as {
		ret = append(ret, assign.Arg{Value: a})
	}
	return
}
