package render

import (
	"git.sr.ht/~ionous/iffy/affine"
	"git.sr.ht/~ionous/iffy/dl/core"
	"git.sr.ht/~ionous/iffy/rt"
	g "git.sr.ht/~ionous/iffy/rt/generic"
)

// expressions are text patterns... so for now adapt via text
// expressions would ideally adapt based on the pattern type
// the assembler probably needs to work directly on tokens...
func (op *RenderPattern) GetText(run rt.Runtime) (ret g.Value, err error) {
	buf := core.Buffer{core.MakeActivity(&core.CallPattern{
		Pattern:   op.Pattern,
		Arguments: op.Arguments,
	})}
	return buf.GetText(run)
}

func (op *RenderPattern) GetAssignedValue(run rt.Runtime) (ret g.Value, err error) {
	det := core.CallPattern{Pattern: op.Pattern, Arguments: op.Arguments}
	if v, e := det.DetermineValue(run); e != nil {
		err = e
	} else {
		ret = v
	}
	return
}

func (op *RenderPattern) Affinity() (ret affine.Affinity) {
	return // fix: we could know in advance... but currently we dont.
}
