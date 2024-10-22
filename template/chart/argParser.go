package chart

import (
	"git.sr.ht/~ionous/tapestry/template/postfix"
)

// ArgParser reads arguments specified by a function call.
type ArgParser struct {
	out     postfix.Expression
	cnt     int
	err     error
	factory ExpressionStateFactory
}

// MakeArgParser using a factory so that tests can mock out recursion;
// Normally, arguments can be any operand -- or any directive.
func MakeArgParser(f ExpressionStateFactory) ArgParser {
	return ArgParser{factory: f}
}

func (p *ArgParser) String() string {
	return "arg parser"
}

// GetArguments returns both the expression, and the number of separated arguments.
func (p *ArgParser) GetArguments() (postfix.Expression, int, error) {
	return p.out, p.cnt, p.err
}

// GetExpression see also GetArguments.
func (p *ArgParser) GetExpression() (postfix.Expression, error) {
	return p.out, p.err
}

// NewRune starts with the first character of an argument;
// each arg is parsed via ArgParser.
func (p *ArgParser) NewRune(r rune) State {
	var sub ExpressionState
	if f := p.factory; f != nil {
		sub = p.factory.NewExpressionState()
	} else {
		sub = new(SubdirParser)
	}
	return RunStep(r, sub, Statement("continuing args", func(r rune) (ret State) {
		if exp, e := sub.GetExpression(); e != nil {
			p.err = e
		} else if len(exp) > 0 {
			// doesn't shunt: operands are already in postfix order.
			p.cnt, p.out = p.cnt+1, append(p.out, exp...)
			if isSpace(r) {
				ret = Step(spaces, p) // loop...
			}
		}
		return
	}))
}
