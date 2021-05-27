package debug

import (
	"git.sr.ht/~ionous/iffy/rt"
)

func (DoNothing) Execute(rt.Runtime) error { return nil }
