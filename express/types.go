package express

import (
	r "reflect"

	"git.sr.ht/~ionous/tapestry/dl/math"
	"git.sr.ht/~ionous/tapestry/rt"
)

var typeNumEval = r.TypeOf((*rt.NumEval)(nil)).Elem()
var typeTextEval = r.TypeOf((*rt.TextEval)(nil)).Elem()
var compareNum = r.TypeOf((*math.CompareNum)(nil)).Elem()
var compareText = r.TypeOf((*math.CompareText)(nil)).Elem()

func implements(a, b r.Value, t r.Type) bool {
	return a.Type().Implements(t) && b.Type().Implements(t)
}
