package game

import (
	"io"

	"git.sr.ht/~ionous/tapestry/rt"
	"github.com/ionous/errutil"
)

// fix: need to think about this. requiring package files or even runtime debug seems iffy.
func (*PrintVersion) Execute(run rt.Runtime) (_ error) {
	//details := false
	//version := files.GetVersion(details)
	io.WriteString(run.Writer(), "???")
	return
}

// returns SignalQuit
// any prompting of the user ( ie. are you sure? ) should happen *before* quit is called.
func (*QuitGame) Execute(run rt.Runtime) error {
	return rt.SignalQuit
}

// todo: it'd be nice for the user to be able to type a name for the file as part of the thing
// including a "request list of save files" type of command
// would either need a parser action that can transparently eat words "load name of my save"
// ( and change Signal into a struct with data )
// -- could also have the play engine handle the player interaction
// -- but that seems less flexible.
func (*LoadGame) Execute(run rt.Runtime) error {
	return rt.SignalLoad
}

func (*SaveGame) Execute(run rt.Runtime) error {
	return rt.SignalSave
}

func (*UndoTurn) Execute(run rt.Runtime) error {
	return errutil.New("undo not implemented")
}
