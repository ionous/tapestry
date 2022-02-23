package core

import (
	"strconv"

	"git.sr.ht/~ionous/tapestry/rt"
)

func MakeActivity(exe ...rt.Execute) []rt.Execute {
	return exe
}

func Args(from ...rt.Assignment) (ret []rt.Arg) {
	for i, from := range from {
		ret = append(ret, rt.Arg{
			Name: W("$" + strconv.Itoa(i+1)),
			From: from,
		})
	}
	return
}

func NamedArgs(name string, from rt.Assignment) []rt.Arg {
	return []rt.Arg{{
		Name: W(name),
		From: from,
	}}
}
