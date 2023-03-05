package story

import (
	"git.sr.ht/~ionous/tapestry/dl/literal"
)

func B(b bool) *literal.BoolValue       { return &literal.BoolValue{Value: b} }
func I(n int) *literal.NumValue         { return &literal.NumValue{Value: float64(n)} }
func F(n float64) *literal.NumValue     { return &literal.NumValue{Value: n} }
func T(s string) *literal.TextValue     { return &literal.TextValue{Value: s} }
func Tx(s, t string) *literal.TextValue { return &literal.TextValue{Value: s, Kind: t} }
