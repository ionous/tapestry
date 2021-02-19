package test

import (
	"git.sr.ht/~ionous/iffy/affine"
	"git.sr.ht/~ionous/iffy/dl/core"
	"git.sr.ht/~ionous/iffy/dl/pattern"
	"git.sr.ht/~ionous/iffy/rt"
	g "git.sr.ht/~ionous/iffy/rt/generic"
)

type MatchGroups struct {
	A       GroupSettings
	B       GroupSettings
	Matches bool `if:"bool"`
}

// a pattern for matching groups --
// we add rules that if things arent equal we return false
var matchGroups = pattern.Pattern{
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
			&core.GetAtField{
				From:  &core.FromVar{N("a")},
				Field: "label",
			},
			&core.NotEqualTo{},
			&core.GetAtField{
				From:  &core.FromVar{N("b")},
				Field: "label",
			},
		},
		Execute: matches(false),
	}, {
		Filter: &core.CompareText{
			&core.GetAtField{
				From:  &core.FromVar{N("a")},
				Field: "innumerable",
			},
			&core.NotEqualTo{},
			&core.GetAtField{
				From:  &core.FromVar{N("b")},
				Field: "innumerable",
			},
		},
		Execute: matches(false),
	}, {
		Filter: &core.CompareText{
			&core.GetAtField{
				From:  &core.FromVar{N("a")},
				Field: "group_options",
			},
			&core.NotEqualTo{},
			&core.GetAtField{
				From:  &core.FromVar{N("b")},
				Field: "group_options",
			},
		},
		Execute: matches(false),
	}},
}

func matches(b bool) rt.Execute {
	return &core.Assign{core.Variable{Str: "matches"}, &core.FromBool{&core.Bool{b}}}
}
