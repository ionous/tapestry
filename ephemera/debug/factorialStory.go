package debug

import (
	"git.sr.ht/~ionous/iffy/dl/core"
	"git.sr.ht/~ionous/iffy/dl/value"
	"git.sr.ht/~ionous/iffy/ephemera/story"
	"git.sr.ht/~ionous/iffy/rt"
)

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
					Value: &core.Activity{[]rt.Execute{
						&core.SayText{
							Text: &core.PrintNum{
								Num: &core.CallPattern{
									Pattern: factorialName,
									Arguments: core.CallArgs{
										Args: []core.CallArg{
											core.CallArg{
												Name: "num",
												From: &core.FromNum{
													Val: &core.NumValue{Num: 3}},
											}},
									}},
							}},
					}},
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
							Value: &core.Activity{[]rt.Execute{
								&core.Assign{
									Var: numVar,
									From: &core.FromNum{&core.ProductOf{
										A: &core.GetVar{Name: numVar},
										B: &core.DiffOf{
											A: &core.GetVar{Name: numVar},
											B: &core.NumValue{Num: 1}},
									}},
								},
							}},
						}},
					}}},
			&story.PatternActions{
				Name:          factorialName,
				PatternReturn: &story.PatternReturn{Result: numberDecl},
				PatternRules: story.PatternRules{
					PatternRule: []story.PatternRule{{
						Guard: &core.CompareNum{
							A:  &core.GetVar{Name: numVar},
							Is: &core.Equal{},
							B:  &core.NumValue{}},
						Hook: story.ProgramHook{
							Choice: story.ProgramHook_Activity_Opt,
							Value: &core.Activity{[]rt.Execute{
								&core.Assign{
									Var:  numVar,
									From: &core.FromNum{&core.NumValue{Num: 1}},
								}},
							}}},
					}},
			}},
	}},
}

var factorialName = value.PatternName{Str: "factorial"}
var numVar = value.VariableName{Str: "num"}

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
