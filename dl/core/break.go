package core

import (
	"git.sr.ht/~ionous/iffy/dl/composer"
	"git.sr.ht/~ionous/iffy/rt"
)

type Break struct{}

type Next struct{}

// DoInterrupt - an error code to break out of loops
type DoInterrupt struct{ KeepGoing bool }

func (e DoInterrupt) Error() string {
	return "loop interrupted"
}

func (*Break) Execute(rt.Runtime) error {
	return DoInterrupt{KeepGoing: false}
}

func (*Next) Execute(rt.Runtime) error {
	return DoInterrupt{KeepGoing: true}
}

func (*Break) Compose() composer.Spec {
	return composer.Spec{
		Fluent: &composer.Fluid{Role: composer.Command},
		Group:  "flow",
		Desc:   "In a repeating loop, exit the loop.",
	}
}

func (*Next) Compose() composer.Spec {
	return composer.Spec{
		Fluent: &composer.Fluid{Role: composer.Command},
		Group:  "flow",
		Desc:   "In a repeating loop, try the next iteration of the loop.",
	}
}
