package debug

import (
	"git.sr.ht/~ionous/tapestry/dl/assign"
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
			&story.DefinePattern{
				PatternName: T("factorial"),
				Params: []story.FieldDefinition{
					&story.NumberField{
						Markup: UserComment("just one argument, a number called 'num'"),
						Name:   "num",
					}},
				Result: &story.NumberField{
					Markup: UserComment("the result uses the same variable as the pattern input does"),
					Name:   "num",
				},
				Rules: []story.PatternRule{{
					Markup: UserComment("rules within a set of rules are evaluated top to bottom"),
					Guard:  FactorialIsZero,
					Does:   FactorialUseOne,
				}, {
					Guard: &core.Always{},
					Does:  FactorialMulMinusOne,
				}},
			}},
	}},
}

// run 3! factorial
var FactorialCheck = []rt.Execute{
	&debug.Expect{
		Value: &core.CompareNum{
			A: F(6), Is: core.Equal, B: &assign.CallPattern{
				PatternName: "factorial",
				Arguments: []assign.Arg{{
					Name:  "num",
					Value: &assign.FromNumber{Value: F(3)},
					// fix: for some reason, the comment isn't appearing in the output.
					// Markup: UserComment("start the factorial with '3'"),
				}}},
		},
	},
}

// subtracts 1 from the num and multiples by one
var FactorialMulMinusOne = []rt.Execute{
	&assign.SetValue{
		Target: assign.Variable("num"),
		Value: &assign.FromNumber{Value: &core.ProductOf{
			A: assign.Variable("num"),
			B: &core.DiffOf{
				A: assign.Variable("num"),
				B: I(1),
			},
		}}},
}

// at 0, use the number 1
var FactorialUseOne = []rt.Execute{
	&assign.SetValue{
		Target: assign.Variable("num"),
		Value:  &assign.FromNumber{Value: I(1)},
	},
}

var FactorialIsZero = &core.CompareNum{
	Markup: UserComment("so, when we've reached 0..."),
	A:      assign.Variable("num"),
	Is:     core.Equal,
	B:      F(0)}
