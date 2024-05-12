package story

import (
	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/rt"
)

// func I(n int) *literal.NumValue         { return &literal.NumValue{Value: float64(n)} }
// func F(n float64) *literal.NumValue     { return &literal.NumValue{Value: n} }
// func T(s string) *literal.TextValue     { return &literal.TextValue{Value: s} }

func truly() rt.Assignment {
	return &assign.FromBool{
		Value: &literal.BoolValue{Value: true},
	}
}
