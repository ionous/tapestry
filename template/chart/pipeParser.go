package chart

import (
	"git.sr.ht/~ionous/tapestry/template/postfix"
	"github.com/ionous/errutil"
)

// Pipe parser reads an expression followed by an optional series of functions.
// Expression | Function | ...
type PipeParser struct {
	err error
	xs  postfix.Expression
}

func (p *PipeParser) String() string {
	return "pipes"
}

// NewRune starts on the first character of an operand or opening sub-phrase.
func (p *PipeParser) NewRune(r rune) State {
	var expParser ExpressionParser
	return RunStep(r, &expParser, p.after(0, &expParser))
}

func (p *PipeParser) GetExpression() (postfix.Expression, error) {
	return p.xs, p.err
}

// after generates a state which reads the results of the passed expression parser.
func (p *PipeParser) after(n int, expParser ExpressionState) State {
	return Statement("check pipe", func(r rune) (ret State) {
		if xs, e := expParser.GetExpression(); e != nil {
			p.err = e
		} else if cnt := len(xs); n > 0 && cnt == n {
			p.err = errutil.New("pipe should be followed by a call")
		} else {
			switch {
			case isPipe(r):
				ret = Statement("post pipe", func(r rune) State {
					// pass the existing expression into the call parser.
					call := CallParser{arity: 1, out: xs}
					return RunStep(r, &call, p.after(cnt, &call))
				})
			default:
				p.xs = xs
			}
		}
		return
	})
}
