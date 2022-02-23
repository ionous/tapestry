package render

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/dl/core"
	"git.sr.ht/~ionous/tapestry/rt"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
)

// expressions are text patterns... so for now adapt via text
// expressions would ideally adapt based on the pattern type
// the assembler probably needs to work directly on tokens...
func (op *RenderPattern) GetText(run rt.Runtime) (ret g.Value, err error) {
	buf := core.BufferText{Does: core.MakeActivity(&op.Call)}
	return buf.GetText(run)
}

func (op *RenderPattern) GetBool(run rt.Runtime) (ret g.Value, err error) {
	return op.Call.GetBool(run)
}

func (op *RenderPattern) GetAssignedValue(run rt.Runtime) (ret g.Value, err error) {
	if v, e := op.Call.DetermineValue(run); e != nil {
		err = e
	} else {
		ret = v
	}
	return
}

func (op *RenderPattern) Affinity() (ret affine.Affinity) {
	return // fix: we could know in advance... but currently we dont.
}
