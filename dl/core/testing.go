package core

import (
	"strconv"

	"git.sr.ht/~ionous/tapestry/rt"
)

func MakeActivity(exe ...rt.Execute) []rt.Execute {
	return exe
}

// takes any of the rt evals
func MakeArgs(as ...Assignment) (ret []Arg) {
	for i, a := range as {
		ret = append(ret, Arg{
			// FIX: this is silly, just have no name and count the args when they are used.
			// in which case MakeArgs itself should be removed.
			Name:  W("$" + strconv.Itoa(i+1)),
			Value: a,
		})
	}
	return
}
