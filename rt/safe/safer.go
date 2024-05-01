package safe

import (
	"fmt"

	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/rt"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
)

func Check(v g.Value, want affine.Affinity) (err error) {
	if have := v.Affinity(); len(want) > 0 && want != have {
		err = fmt.Errorf("wanted %q, have %q", want, have)
	}
	return
}

func CheckList(v g.Value) (err error) {
	if have := v.Affinity(); !affine.IsList(have) {
		err = fmt.Errorf("wanted a list, have %q", have)
	}
	return
}

// fix! ( at the very least should live in pattern
// but we need to remove its few -- tests and Determine -- dependencies on core
var HackTillTemplatesCanEvaluatePatternTypes g.Value

func GetTemplateText() (ret g.Value) {
	if hack := HackTillTemplatesCanEvaluatePatternTypes; hack != nil {
		// we didn't accumulate any text during execution
		// but perhaps we ran a pattern that returned text.
		// to get rid of this, we'd examine (at runtime or compile time) the futures calls
		// and switch on execute patterns vs text patterns
		// an example is { .Lantern } which says the name
		// vs. { pluralize: .Lantern } which returns the pluralized name.
		ret = hack
		HackTillTemplatesCanEvaluatePatternTypes = nil
	} else {
		ret = g.Empty // if the res was empty, it might have intentionally been empty
	}
	return
}

func PopSeveral(run rt.Runtime, p int) {
	for i := 0; i < p; i++ {
		run.PopScope()
	}
}
