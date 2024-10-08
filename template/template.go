package template

import (
	"git.sr.ht/~ionous/tapestry/template/chart"
	"git.sr.ht/~ionous/tapestry/template/postfix"
	"github.com/ionous/errutil"
)

// Expression provides a local alias for postfix.Expression;
// a series of postfix.Function records.
type Expression = postfix.Expression

// Parse the passed template string into an expression.
func Parse(template string) (ret Expression, err error) {
	p := chart.MakeTemplateParser()
	e := chart.Parse(template, &p)
	xs, ex := p.GetExpression()
	if ex != nil {
		err = errutil.New(ex, e)
	} else if e != nil {
		err = e
	} else {
		ret = xs
	}
	return
}

// ParseExpression reads a series of simple operand and operator phrases
// and creates a series of postfix.Function records.
// ex. "(5+6)*(1+2)" -> 5 6 ADD 1 2 ADD MUL
// where MUL and ADD are types.Operator,
// while the numbers are types.Number.
func ParseExpression(str string) (ret Expression, err error) {
	var p chart.SeriesParser
	if e := chart.Parse(str, &p); e != nil {
		err = e
	} else {
		ret, err = p.GetExpression()
	}
	return
}
