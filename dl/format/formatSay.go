package format

import (
	"io"

	"git.sr.ht/~ionous/tapestry/dl/cmd"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/safe"
)

func (op *SayActor) Execute(run rt.Runtime) (err error) {
	if _, e := safe.GetText(run, op.Actor); e != nil {
		err = cmd.Error(op, e)
	} else if text, e := safe.GetText(run, op.Text); e != nil {
		err = cmd.Error(op, e)
	} else {
		io.WriteString(run.Writer(), text.String())
	}
	return
}
