package core

import (
	"git.sr.ht/~ionous/iffy/affine"
	"git.sr.ht/~ionous/iffy/dl/composer"
	"git.sr.ht/~ionous/iffy/rt"

	g "git.sr.ht/~ionous/iffy/rt/generic"
	"git.sr.ht/~ionous/iffy/rt/safe"
)

// Send triggers an event. All events are expected to return a "bool" success.
// erroring if the value requested doesnt support the error returned.
type Send struct {
	Event     EventName       // a text eval here would be like a function pointer maybe...
	Path      rt.TextListEval // object names
	Arguments *Arguments      // event args shouldnt be optional, but it mirrors pattern
}

type EventName string

func (n EventName) String() string { return string(n) }

func (*Send) Compose() composer.Spec {
	return composer.Spec{
		Spec:  "Send: {event:event_name} {?arguments} to:{path:text_list_eval}",
		Group: "events",
		Desc:  "Send: Triggers a event, returns a true/false success value.",
	}
}

func (op *Send) Execute(run rt.Runtime) error {
	_, err := op.send(run, "")
	return err
}

// GetBool returns the first matching bool evaluation.
func (op *Send) GetBool(run rt.Runtime) (g.Value, error) {
	return op.send(run, affine.Bool)
}

func (op *Send) send(run rt.Runtime, aff affine.Affinity) (ret g.Value, err error) {
	if path, e := safe.GetTextList(run, op.Path); e != nil {
		err = e
	} else {
		var args []rt.Arg
		for _, a := range op.Arguments.Args {
			args = append(args, rt.Arg{a.Name, a.From})
		}
		name, up := op.Event.String(), path.Strings()
		if v, e := run.Send(name, up, args); e != nil {
			err = cmdErrorCtx(op, name, e)
		} else {
			ret = v
		}
	}
	return
}
