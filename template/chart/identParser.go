package chart

import "strings"

type IdentParser struct {
	runes strings.Builder
}

func (p *IdentParser) String() string {
	return "identifiers"
}

func (p *IdentParser) Reset() string {
	r := p.runes.String()
	p.runes.Reset()
	return r
}

func (p *IdentParser) Identifier() string {
	return p.runes.String()
}

// first character of the identifier
func (p *IdentParser) NewRune(r rune) (ret State) {
	if isLetter(r) {
		p.runes.WriteRune(r)
		ret = Statement("ident head", p.body) // loop...
	}
	return
}

// subsequent characters can be letters or numbers
// noting that fields are separated by dots "."
func (p *IdentParser) body(r rune) (ret State) {
	if isLetter(r) || isNumber(r) || isQualifier(r) {
		p.runes.WriteRune(r)
		ret = Statement("ident body", p.body) // loop...
	}
	return
}
