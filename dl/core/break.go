package core

import (
	"git.sr.ht/~ionous/iffy/rt"
)

// DoInterrupt - an error code to break out of loops
type DoInterrupt struct{ KeepGoing bool }

func (e DoInterrupt) Error() string {
	return "loop interrupted"
}

func (e DoInterrupt) NoPanic() {}

func (*Break) Execute(rt.Runtime) error {
	return DoInterrupt{KeepGoing: false}
}

func (*Next) Execute(rt.Runtime) error {
	return DoInterrupt{KeepGoing: true}
}
