package debug

import (
	"git.sr.ht/~ionous/tapestry/dl/core"
	"git.sr.ht/~ionous/tapestry/dl/debug"
	"git.sr.ht/~ionous/tapestry/dl/story"
	"git.sr.ht/~ionous/tapestry/rt"
)

func UserComment(s string) map[string]any {
	return map[string]any{"comment": s}
}

// a program that can check factorials
var FactorialStory = &story.Story{
	Paragraph: []story.Paragraph{{
		StoryStatement: []story.StoryStatement{
			&debug.Test{
				TestName: debug.TestName{
					Str: "factorial",
				},
				Do: FactorialCheck,
			},
			&story.PatternDecl{
				Name: factorialName,
				PatternReturn: &story.PatternReturn{Result: &story.NumberField{
					Markup: UserComment("the result uses the same variable as the pattern input does"),
					Name:   "num",
				}},
				Params: []story.Field{
					&story.NumberField{
						Markup: UserComment("just one argument, a number called 'num'"),
						Name:   "num",
					}},
			},
			&story.PatternActions{
				Name: factorialName,
				Rules: []story.PatternRule{{
					Guard: &core.Always{},
					Does:  FactorialMulMinusOne,
				}}},
			&story.PatternActions{
				Name: factorialName,
				Rules: []story.PatternRule{{
					Markup: UserComment("the rule considered first is the rule that was written last:"),
					Guard:  FactorialIsZero,
					Does:   FactorialUseOne,
				}},
			}},
	}},
}

// run 3! factorial
var FactorialCheck = []rt.Execute{
	&debug.Expect{
		Value: &core.CompareNum{
			A: F(6), Is: core.Equal, B: &core.CallPattern{
				Pattern: factorialName,
				Arguments: []core.Arg{core.Arg{
					Name:  "num",
					Value: core.AssignFromNumber(F(3)),
					// fix: for some reason, the comment isn't appearing in the output.
					Markup: UserComment("start the factorial with '3'"),
				}}},
		},
	},
}

// subtracts 1 from the num and multiples by one
var FactorialMulMinusOne = []rt.Execute{
	&core.SetValue{
		Target: core.Variable("num"),
		Value: core.AssignFromNumber(&core.ProductOf{
			A: GetVariable("num"),
			B: &core.DiffOf{
				A: GetVariable("num"),
				B: I(1),
			},
		})},
}

// at 0, use the number 1
var FactorialUseOne = []rt.Execute{
	&core.SetValue{
		Target: core.Variable("num"),
		Value:  core.AssignFromNumber(I(1)),
	},
}

var FactorialIsZero = &core.CompareNum{
	Markup: UserComment("so, when we've reached 0..."),
	A:      GetVariable("num"),
	Is:     core.Equal,
	B:      F(0)}

var factorialName = core.PatternName{Str: "factorial"}
