package test

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/dl/core"
	"git.sr.ht/~ionous/tapestry/rt"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"git.sr.ht/~ionous/tapestry/test/testpat"
)

type MatchGroups struct {
	A       GroupSettings
	B       GroupSettings
	Matches bool `if:"bool"`
}

// a pattern for matching groups --
// we add rules that if things arent equal we return false
var matchGroups = testpat.Pattern{
	Name:   "match_groups",
	Labels: []string{"a", "b"},
	Fields: []g.Field{
		{Name: "a", Affinity: affine.Record, Type: "group_settings"},
		{Name: "b", Affinity: affine.Record, Type: "group_settings"},
		{Name: "matches", Affinity: affine.Bool},
	},
	Return: "matches",
	// rules are evaluated in reverse order ( see sortRules )
	Rules: []rt.Rule{{
		Filter:  &core.Always{},
		Execute: matches(true),
	}, {
		Filter: &core.CompareText{
			A: &core.GetAtField{
				From:  &core.FromVar{Var: N("a")},
				Field: W("label"),
			},
			Is: core.Unequal,
			B: &core.GetAtField{
				From:  &core.FromVar{Var: N("b")},
				Field: W("label"),
			},
		},
		Execute: matches(false),
	}, {
		Filter: &core.CompareText{
			A: &core.GetAtField{
				From:  &core.FromVar{Var: N("a")},
				Field: W("innumerable"),
			},
			Is: core.Unequal,
			B: &core.GetAtField{
				From:  &core.FromVar{Var: N("b")},
				Field: W("innumerable"),
			},
		},
		Execute: matches(false),
	}, {
		Filter: &core.CompareText{
			A: &core.GetAtField{
				From:  &core.FromVar{Var: N("a")},
				Field: W("group_options"),
			},
			Is: core.Unequal,
			B: &core.GetAtField{
				From:  &core.FromVar{Var: N("b")},
				Field: W("group_options"),
			},
		},
		Execute: matches(false),
	}},
}

func matches(b bool) []rt.Execute {
	return []rt.Execute{&core.Assign{Var: N("matches"), From: &core.FromBool{Val: B(b)}}}
}
