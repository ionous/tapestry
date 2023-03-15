package test

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/dl/assign"
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
			A:  core.Variable("a", "label"),
			Is: core.Unequal,
			B:  core.Variable("b", "label"),
		},
		Execute: matches(false),
	}, {
		Filter: &core.CompareText{
			A:  core.Variable("a", "innumerable"),
			Is: core.Unequal,
			B:  core.Variable("b", "innumerable"),
		},
		Execute: matches(false),
	}, {
		Filter: &core.CompareText{
			A:  core.Variable("a", "group_options"),
			Is: core.Unequal,
			B:  core.Variable("b", "group_options"),
		},
		Execute: matches(false),
	}},
}

func matches(b bool) []rt.Execute {
	return []rt.Execute{
		&assign.SetValue{
			Target: core.Variable("matches"),
			Value:  &assign.FromBool{Value: B(b)}},
	}
}
