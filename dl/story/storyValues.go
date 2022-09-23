package story

import (
	"git.sr.ht/~ionous/tapestry/dl/composer"
	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/imp"
	"github.com/ionous/errutil"
)

func B(b bool) *literal.BoolValue       { return &literal.BoolValue{Value: b} }
func I(n int) *literal.NumValue         { return &literal.NumValue{Value: float64(n)} }
func F(n float64) *literal.NumValue     { return &literal.NumValue{Value: n} }
func T(s string) *literal.TextValue     { return &literal.TextValue{Value: s} }
func Tx(s, t string) *literal.TextValue { return &literal.TextValue{Value: s, Kind: t} }

func (op *Certainty) ImportString(k *imp.Importer) (ret string, err error) {
	if str, ok := composer.FindChoice(op, op.Str); !ok {
		err = ImportError(op, errutil.Fmt("%w %q", InvalidValue, op.Str))
	} else {
		ret = str
	}
	return
}
