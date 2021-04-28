package render

import (
	"git.sr.ht/~ionous/iffy/affine"
	"git.sr.ht/~ionous/iffy/dl/composer"
	"git.sr.ht/~ionous/iffy/dl/core"
	"git.sr.ht/~ionous/iffy/rt"
	g "git.sr.ht/~ionous/iffy/rt/generic"
)

// printing is generally an activity b/c say is an activity
// and we want the ability to say several things in series.
// expressions are text patterns... so for now adapt via text
// expressions would ideally adapt based on the pattern type
// the assembler probably needs to work directly on tokens...
type RenderPattern struct {
	core.Determine
}

func (op *RenderPattern) Compose() composer.Spec {
	return composer.Spec{
		Lede:  "render",
		Group: "internal",
	}
}

func (op *RenderPattern) GetText(run rt.Runtime) (ret g.Value, err error) {
	buf := core.Buffer{core.MakeActivity(&core.Determine{
		Pattern:   op.Pattern,
		Arguments: op.Arguments,
	})}
	return buf.GetText(run)
}

func (op *RenderPattern) GetAssignedValue(run rt.Runtime) (ret g.Value, err error) {
	if v, e := op.Determine.DetermineValue(run); e != nil {
		err = e
	} else {
		ret = v
	}
	return
}

func (op *RenderPattern) Affinity() (ret affine.Affinity) {
	return // fix: we could know in advance... but currently we dont.
}
