package express

import (
	r "reflect"

	"git.sr.ht/~ionous/tapestry/dl/core"
	"git.sr.ht/~ionous/tapestry/rt"
)

var typeNumEval = r.TypeOf((*rt.NumEval)(nil)).Elem()
var typeTextEval = r.TypeOf((*rt.TextEval)(nil)).Elem()
var compareNum = r.TypeOf((*core.CompareNum)(nil)).Elem()
var compareText = r.TypeOf((*core.CompareText)(nil)).Elem()

func implements(a, b r.Value, t r.Type) bool {
	return a.Type().Implements(t) && b.Type().Implements(t)
}
