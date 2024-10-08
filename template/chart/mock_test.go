package chart

import (
	"strings"
	"unicode"

	"git.sr.ht/~ionous/tapestry/template/postfix"
	"git.sr.ht/~ionous/tapestry/template/types"
)

// EmptyFactory creates parsers which match empty test.
// implements ExpressionStateFactory
type EmptyFactory struct{}

// EmptyParser reads empty text.
type EmptyParser struct{}

// AnyFactory creates parsers which match any series of lowercase letters
// implements ExpressionStateFactory
type AnyFactory struct{}

// AnyParser reads letters.
type AnyParser struct{ runes strings.Builder }

func (*EmptyParser) String() string {
	return "empty"
}

func (*EmptyFactory) NewExpressionState() ExpressionState {
	return &EmptyParser{}
}
func (*EmptyParser) GetExpression() (ret postfix.Expression, err error) {
	return
}
func (p *EmptyParser) NewRune(r rune) (ret State) {
	if isSpace(r) {
		ret = p
	}
	return
}

func (f *AnyFactory) NewExpressionState() ExpressionState {
	return &AnyParser{}
}

func (*AnyParser) String() string {
	return "any"
}

func (p *AnyParser) GetExpression() (ret postfix.Expression, err error) {
	if s := p.runes.String(); len(s) > 0 {
		arg := types.Reference([]string{s})
		ret = append(ret, arg)
	}
	return
}

func (p *AnyParser) NewRune(r rune) (ret State) {
	if unicode.IsLower(r) {
		p.runes.WriteRune(r)
		ret = p
	} else if isDot(r) {
		ret = p
	}
	return
}
