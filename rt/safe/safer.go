package safe

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/rt"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"github.com/ionous/errutil"
)

func Check(v g.Value, want affine.Affinity) (err error) {
	if va := v.Affinity(); len(want) > 0 && want != va {
		err = errutil.Fmt("wanted %q, have %q", want, va)
	}
	return
}

// fix! ( at the very least should live in pattern
// but we need to remove its few -- tests and Determine -- dependencies on core
var HackTillTemplatesCanEvaluatePatternTypes g.Value

func PopSeveral(run rt.Runtime, p int) {
	for i := 0; i < p; i++ {
		run.PopScope()
	}
}
