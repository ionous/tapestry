package debug

import (
	"git.sr.ht/~ionous/tapestry/dl/core"
	"git.sr.ht/~ionous/tapestry/dl/prim"
	"git.sr.ht/~ionous/tapestry/dl/story"
	"git.sr.ht/~ionous/tapestry/rt"
)

// a program that can check factorials
var FactorialStory = &story.Story{
	Paragraph: []story.Paragraph{{
		StoryStatement: []story.StoryStatement{
			&story.TestStatement{
				TestName: story.TestName{
					Str: "factorial",
				},
				Test: &story.TestOutput{
					UserComment: "3! should equal 6",
					Lines: prim.Lines{
						Str: "6",
					}}},
			&story.TestRule{
				TestName: story.TestName{
					Str: "factorial",
				},
				Does: FactorialCheck,
			},
			&story.PatternDecl{
				Name: factorialName,
				PatternReturn: &story.PatternReturn{Result: &story.NumberField{
					UserComment: "the result uses the same variable as the pattern input does",
					Name:        numVar.Str,
				}},
				Params: []story.Field{
					&story.NumberField{
						UserComment: "just one argument, a number called 'num'",
						Name:        numVar.Str,
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
					UserComment: "the rule considered first is the rule that was written last:",
					Guard:       FactorialIsZero,
					Does:        FactorialUseOne,
				}},
			}},
	}},
}

// run 3! factorial
var FactorialCheck = []rt.Execute{
	&core.SayText{
		Text: &core.PrintNum{
			Num: &core.CallPattern{
				Pattern: factorialName,
				Arguments: []rt.Arg{rt.Arg{
					Name: "num",
					From: &core.FromNum{
						UserComment: "start the factorial with '3'",
						Val:         F(3),
					},
				}}},
		}},
}

// subtracts 1 from the num and multiples by one
var FactorialMulMinusOne = []rt.Execute{
	&core.Assign{
		Var: numVar,
		From: &core.FromNum{Val: &core.ProductOf{
			A: &core.GetVar{Name: numVar},
			B: &core.DiffOf{
				A: &core.GetVar{Name: numVar},
				B: F(1)},
		}},
	},
}

// at 0, use the number 1
var FactorialUseOne = []rt.Execute{
	&core.Assign{
		Var: numVar,
		From: &core.FromNum{
			Val:         F(1),
			UserComment: "...return 1.",
		},
	},
}

var FactorialIsZero = &core.CompareNum{
	UserComment: "so, when we've reached 0...",
	A:           &core.GetVar{Name: numVar},
	Is:          &core.Equal{},
	B:           F(0)}

var factorialName = core.PatternName{Str: "factorial"}
var numVar = core.VariableName{Str: "num"}
