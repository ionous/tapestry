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

// placeholder for ending the game.
// not 100% clear if this will ultimately be a script action, or just launch script actions:
// traditionally prints `*** The End ***`, with an optionally customizable message.
// shifts the parser state ( maybe based on the game state )
// and a prompt containing the parser state options.
func (*EndGame) Execute(run rt.Runtime) error {
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
