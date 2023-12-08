package chart

import "strings"

// KeyParser reads a key and its optional following expression.
type KeyParser struct {
	runes strings.Builder
	exp   ExpressionState
}

func (p *KeyParser) String() string {
	return "keys"
}

// NewRune starts on the first letter of the key.
func (p *KeyParser) NewRune(r rune) (ret State) {
	if isLetter(r) {
		p.runes.WriteRune(r)
		ret = p
	} else if isSpace(r) {
		ret = Step(spaces, p.exp)
	}
	return
}

func (p *KeyParser) GetDirective() (ret Directive, err error) {
	if exp, e := p.exp.GetExpression(); e != nil {
		err = e
	} else {
		ret.Key = p.runes.String()
		ret.Expression = exp
	}
	return
}
