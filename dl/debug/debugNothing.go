package debug

import (
	"git.sr.ht/~ionous/tapestry/rt"
)

type DoNothingText string

func (DoNothing) Execute(rt.Runtime) error { return nil }
