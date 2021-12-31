package chart

import (
	"git.sr.ht/~ionous/tapestry/template/postfix"
	"git.sr.ht/~ionous/tapestry/template/types"
)

// implements OperandState.
type BooleanParser struct {
	IdentParser
}

func (p *BooleanParser) StateName() string {
	return "bools"
}

func (p *BooleanParser) GetOperand() (ret postfix.Function, err error) {
	switch id := p.IdentParser.Identifier(); id {
	case "true":
		ret = types.Bool(true)
	case "false":
		ret = types.Bool(false)
	}
	return
}
