package core

import (
	"git.sr.ht/~ionous/iffy/dl/composer"
	"git.sr.ht/~ionous/iffy/object"
	"git.sr.ht/~ionous/iffy/rt"
	g "git.sr.ht/~ionous/iffy/rt/generic"
	"git.sr.ht/~ionous/iffy/rt/safe"
)

// if the option PrintResponseNames is enabled, the
type Response struct {
	Name string      `if:"selector"`
	Text rt.TextEval `if:"optional,selector=with"`
}

func (*Response) Compose() composer.Spec {
	return composer.Spec{
		Group:  "output",
		Desc:   "Response: Generate text in a replaceable manner.",
		Fluent: &composer.Fluid{Role: composer.Function},
	}
}

func (op *Response) GetText(run rt.Runtime) (ret g.Value, err error) {
	if safe.GetFlag(run, object.PrintResponseNames) {
		ret = g.StringOf(op.Name)
	} else if v, e := safe.GetText(run, op.Text); e != nil {
		err = cmdError(op, e)
	} else {
		ret = v
	}
	return
}
