package chart

import "git.sr.ht/~ionous/tapestry/template/postfix"

// implements OperandState by returning nothing
type EmptyOperand struct{ r rune }

func (p *EmptyOperand) String() string {
	return "empty operand"
}

func (p *EmptyOperand) NewRune(r rune) State {
	return nil
}

func (p *EmptyOperand) GetOperand() (postfix.Function, error) {
	return nil, nil
}
