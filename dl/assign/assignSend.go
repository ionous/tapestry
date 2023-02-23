package assign

import (
	"git.sr.ht/~ionous/tapestry/rt"
	"github.com/ionous/errutil"

	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"git.sr.ht/~ionous/tapestry/rt/safe"
)

// Execute runs send without returning a value
func (op *CallSend) Execute(run rt.Runtime) (err error) {
	_, err = op.GetBool(run)
	return
}

// GetBool returns the first matching bool evaluation.
func (op *CallSend) GetBool(run rt.Runtime) (ret g.Value, err error) {
	if path, e := safe.GetTextList(run, op.Path); e != nil {
		err = cmdError(op, e)
	} else if evt, ok := op.Event.(*CallPattern); !ok {
		err = cmdError(op, errutil.New("expected call pattern in send"))
	} else {
		name, up := evt.PatternName, path.Strings()
		if rec, e := MakeRecord(run, name, evt.Arguments...); e != nil {
			err = cmdError(op, e)
		} else if v, e := run.Send(rec, up); e != nil {
			err = cmdErrorCtx(op, name, e)
		} else {
			ret = v
		}
	}
	return
}
