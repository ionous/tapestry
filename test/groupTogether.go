package test

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/dl/core"
	"git.sr.ht/~ionous/tapestry/dl/list"
	"git.sr.ht/~ionous/tapestry/rt"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"git.sr.ht/~ionous/tapestry/test/testpat"
)

var runGroupTogther = list.ListMap{
	Target:      assign.Variable("settings"),
	List:        &assign.FromTextList{Value: assign.Variable("objects")},
	PatternName: P("assign_grouping"),
}

type AssignGrouping struct {
	In  string
	Out GroupSettings
}

// from a list of object names, build a list of group settings
var assignGrouping = testpat.Pattern{
	Name:   "assign_grouping",
	Return: "out",
	Labels: []string{"in"},
	Fields: []g.Field{
		{Name: "in", Affinity: affine.Text},
		{Name: "out", Affinity: affine.Record, Type: "group_settings"},
	},
	Rules: []rt.Rule{
		{Execute: []rt.Execute{
			&assign.SetValue{
				Target: assign.Variable("out", "name"),
				Value:  &assign.FromText{Value: assign.Variable("in")}},
			&core.ChooseAction{
				If: &core.Matches{
					Text:  assign.Variable("in"),
					Match: "^thing"},
				Does: core.MakeActivity(
					&assign.SetValue{
						Target: assign.Variable("out", "label"),
						Value:  &assign.FromText{Value: T("thingies")}},
				),
			},
		}},
	},
}
