package debug

import (
	"git.sr.ht/~ionous/tapestry/rt"
)

func (DoNothing) Execute(rt.Runtime) error { return nil }
