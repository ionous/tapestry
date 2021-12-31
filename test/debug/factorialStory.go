package debug

import (
	"git.sr.ht/~ionous/tapestry/dl/core"
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
					Lines: story.Lines{
						Str: "6",
					}}},
			&story.TestRule{
				TestName: story.TestName{
					Str: "factorial",
				},
				Hook: story.ProgramHook{
					Choice: story.ProgramHook_Activity_Opt,
					Value:  FactorialCheck,
				}},
			&story.PatternDecl{
				Name: factorialName,
				Type: story.PatternType{
					Str: story.PatternType_Patterns},
				Optvars: &story.PatternVariablesTail{
					VariableDecl: []story.VariableDecl{numberDecl}},
			},
			&story.PatternActions{
				Name:          factorialName,
				PatternReturn: &story.PatternReturn{Result: numberDecl},
				PatternRules: story.PatternRules{
					PatternRule: []story.PatternRule{{
						Guard: &core.Always{},
						Hook: story.ProgramHook{
							Choice: story.ProgramHook_Activity_Opt,
							Value:  FactorialMulMinusOne,
						}}}}},
			&story.PatternActions{
				Name:          factorialName,
				PatternReturn: &story.PatternReturn{Result: numberDecl},
				PatternRules: story.PatternRules{
					PatternRule: []story.PatternRule{{
						Guard: FactorialIsZero,
						Hook: story.ProgramHook{
							Choice: story.ProgramHook_Activity_Opt,
							Value:  FactorialUseOne,
						}}}},
			}},
	}},
}

// run 3! factorial
var FactorialCheck = &core.Activity{[]rt.Execute{
	&core.SayText{
		Text: &core.PrintNum{
			Num: &core.CallPattern{
				Pattern: factorialName,
				Arguments: core.CallArgs{
					Args: []core.CallArg{
						core.CallArg{
							Name: "num",
							From: &core.FromNum{
								Val: F(3)},
						}},
				}},
		}},
}}

// subtracts 1 from the num and multiples by one
var FactorialMulMinusOne = &core.Activity{[]rt.Execute{
	&core.Assign{
		Var: numVar,
		From: &core.FromNum{&core.ProductOf{
			A: &core.GetVar{Name: numVar},
			B: &core.DiffOf{
				A: &core.GetVar{Name: numVar},
				B: F(1)},
		}},
	},
}}

// at 0, use the number 1
var FactorialUseOne = &core.Activity{[]rt.Execute{
	&core.Assign{
		Var:  numVar,
		From: &core.FromNum{F(1)},
	}},
}

var FactorialIsZero = &core.CompareNum{
	A:  &core.GetVar{Name: numVar},
	Is: &core.Equal{},
	B:  F(0)}

var factorialName = core.PatternName{Str: "factorial"}
var numVar = core.VariableName{Str: "num"}

var numberDecl = story.VariableDecl{
	An: story.Determiner{
		Str: story.Determiner_A,
	},
	Name: numVar,
	Type: story.VariableType{
		Choice: story.VariableType_Primitive_Opt,
		Value: &story.PrimitiveType{
			Str: story.PrimitiveType_Number,
		}},
}
