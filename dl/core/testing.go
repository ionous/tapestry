package core

import (
	"strconv"

	"git.sr.ht/~ionous/tapestry/rt"
)

func MakeActivity(exe ...rt.Execute) []rt.Execute {
	return exe
}

func Args(from ...rt.Assignment) CallArgs {
	var p CallArgs
	for i, from := range from {
		p.Args = append(p.Args, rt.Arg{
			Name: W("$" + strconv.Itoa(i+1)),
			From: from,
		})
	}
	return p
}

func NamedArgs(name string, from rt.Assignment) CallArgs {
	return CallArgs{Args: []rt.Arg{{
		Name: W(name),
		From: from,
	}}}
}
