package debug

import (
	"git.sr.ht/~ionous/iffy/dl/composer"
	"git.sr.ht/~ionous/iffy/rt"
)

// DoNothing implements Execute, but .... does nothing.
type DoNothing struct {
	Reason string
}

func (*DoNothing) Compose() composer.Spec {
	return composer.Spec{
		Group: "exec",
		Desc:  "Do nothing: Statement which does nothing.",
	}
}

func (DoNothing) Execute(rt.Runtime) error { return nil }
