package core

import (
	"git.sr.ht/~ionous/iffy/dl/composer"
	"git.sr.ht/~ionous/iffy/object"
	"git.sr.ht/~ionous/iffy/rt"
	g "git.sr.ht/~ionous/iffy/rt/generic"
)

type HasDominion struct {
	Name string
}

func (*HasDominion) Compose() composer.Spec {
	return composer.Spec{
		Group: "logic",
	}
}

func (op *HasDominion) GetBool(run rt.Runtime) (ret g.Value, err error) {
	return run.GetField(object.Domain, op.Name)
}
