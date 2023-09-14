package sys

import (
	"git.sr.ht/~ionous/tapestry/rt"
	"github.com/ionous/errutil"
)

func (*PrintVersion) Execute(run rt.Runtime) (err error) {
	return errutil.New("version not implemented")
}
func (*QuitGame) Execute(run rt.Runtime) (err error) {
	return errutil.New("quit not implemented")
}
func (*RestoreGame) Execute(run rt.Runtime) (err error) {
	return errutil.New("restore not implemented")
}
func (*SaveGame) Execute(run rt.Runtime) (err error) {
	return errutil.New("save not implemented")
}
func (*UndoTurn) Execute(run rt.Runtime) (err error) {
	return errutil.New("undo not implemented")
}
