package chart

import (
	"git.sr.ht/~ionous/tapestry/template/postfix"
	"github.com/ionous/errutil"
)

// ExpressionParser reads either a single call or a series of operations.
type ExpressionParser struct {
	out        postfix.Expression
	err        error
	argFactory ExpressionStateFactory // for testing
}

func (p *ExpressionParser) StateName() string {
	return "expression"
}
func (p *ExpressionParser) GetExpression() (postfix.Expression, error) {
	return p.out, p.err
}

// NewRune starts with the first character of a directive's content.
func (p *ExpressionParser) NewRune(r rune) State {
	call := MakeCallParser(0, p.argFactory)
	series := SeriesParser{}
	para := Parallel("call or series",
		Step(&series, OnExit("series", func() {
			if x, e := series.GetExpression(); e != nil {
				p.err = errutil.Append(p.err, e)
			} else if len(x) > 0 {
				p.out, p.err = x, nil // longest match wins; noting that operators can read too far ahead
				// particularly an issue with ! which provisionally matches !=
			}
		})),
		Step(&call, OnExit("call", func() {
			if x, e := call.GetExpression(); e != nil {
				p.err = errutil.Append(p.err, e)
			} else if len(x) > 0 {
				p.out, p.err = x, nil // for equal matches, call wins.
			}
		})),
	)
	return para.NewRune(r)
}
