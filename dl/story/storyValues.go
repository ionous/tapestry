package story

import (
	"git.sr.ht/~ionous/tapestry/dl/composer"
	"git.sr.ht/~ionous/tapestry/dl/literal"
	"github.com/ionous/errutil"
)

func B(b bool) *literal.BoolValue   { return &literal.BoolValue{Bool: b} }
func I(n int) *literal.NumValue     { return &literal.NumValue{Num: float64(n)} }
func F(n float64) *literal.NumValue { return &literal.NumValue{Num: n} }
func T(s string) *literal.TextValue { return &literal.TextValue{Text: s} }

func (op *Certainty) ImportString(k *Importer) (ret string, err error) {
	if str, ok := composer.FindChoice(op, op.Str); !ok {
		err = ImportError(op, op.At, errutil.Fmt("%w %q", InvalidValue, op.Str))
	} else {
		ret = str
	}
	return
}
