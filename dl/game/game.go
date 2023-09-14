package game

import (
	"git.sr.ht/~ionous/tapestry/rt"
	"github.com/ionous/errutil"
)

// error code for system commands
// used by play/step to switch take some action as a side-effect a command
type Signal int

//go:generate stringer -type=Signal
const (
	SignalUnknown Signal = iota
	SignalQuit
)

func (s Signal) Error() string {
	return s.String()
}

func (*PrintVersion) Execute(run rt.Runtime) error {
	return errutil.New("version not implemented")
}

// returns SignalQuit
// any prompting of the user ( ie. are you sure? ) should happen *before* quit is called.
func (*QuitGame) Execute(run rt.Runtime) error {
	return SignalQuit
}

func (*RestoreGame) Execute(run rt.Runtime) error {
	return errutil.New("restore not implemented")
}

func (*SaveGame) Execute(run rt.Runtime) error {
	return errutil.New("save not implemented")
}

func (*UndoTurn) Execute(run rt.Runtime) error {
	return errutil.New("undo not implemented")
}
