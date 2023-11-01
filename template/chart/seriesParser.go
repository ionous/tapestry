package chart

import (
	"git.sr.ht/~ionous/tapestry/template/postfix"
	"git.sr.ht/~ionous/tapestry/template/types"
)

// SeriesParser reads a sequence of operand and operator phrases.
type SeriesParser struct {
	err error
	out postfix.Shunt
}

func (p *SeriesParser) String() string {
	return "series"
}

// NewRune starts on the first character of an operand or opening sub-phrase.
func (p *SeriesParser) NewRune(r rune) State {
	return p.operand(r)
}

func (p *SeriesParser) GetExpression() (ret postfix.Expression, err error) {
	if e := p.err; e != nil {
		err = e
	} else {
		ret, err = p.out.GetExpression()
	}
	return
}

// at the start of every operand we might have some opening paren, a bracket,
// or just some operand.
func (p *SeriesParser) operand(r rune) (ret State) {
	switch {
	case isOpenParen(r):
		p.out.BeginSubExpression()
		ret = Step(spaces, Statement("sub expression", p.operand))
	default:
		var a SubdirParser
		ret = RunStep(r, &a, Statement("after sub series", func(r rune) (ret State) {
			if exp, e := a.GetExpression(); e != nil {
				p.err = e
			} else if len(exp) > 0 {
				p.out.AddExpression(exp)
				ret = RunStep(r, spaces, Statement("operator", p.operator))
			}
			return
		}))
	}
	return
}

// after every argument can come operators or close parens or the end
// start on the first character of the operator.
// a pipe floats upward.
func (p *SeriesParser) operator(r rune) (ret State) {
	var b OperatorParser
	return RunStep(r, &b, Statement("after series op", func(r rune) (ret State) {
		switch {
		case isCloseParen(r):
			p.out.EndSubExpression()
			ret = Step(spaces, Statement("closing", p.operator))
		default:
			ret = RunStep(r, &b, Statement("continuing series", func(r rune) (ret State) {
				if op, e := b.GetOperator(); e != nil {
					p.err = e
				} else if op != types.Operator(-1) {
					p.out.AddFunction(op)
					ret = RunStep(r, spaces, Statement("next operand", p.operand))
				}
				return
			}))
		}
		return
	}))
}
