package core

import (
	"strconv"

	"git.sr.ht/~ionous/iffy/rt"
)

func NewActivity(exe ...rt.Execute) *Activity {
	return &Activity{Exe: exe}
}

func MakeActivity(exe ...rt.Execute) Activity {
	return Activity{Exe: exe}
}

func Args(from ...rt.Assignment) CallArgs {
	var p CallArgs
	for i, from := range from {
		p.Args = append(p.Args, CallArg{
			Name: W("$" + strconv.Itoa(i+1)),
			From: from,
		})
	}
	return p
}

func NamedArgs(name string, from rt.Assignment) CallArgs {
	return CallArgs{Args: []CallArg{{
		Name: W(name),
		From: from,
	}}}
}
