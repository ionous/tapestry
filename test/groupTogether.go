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
	Target:      core.Variable("settings"),
	List:        &assign.FromTextList{Value: core.Variable("objects")},
	PatternName: P("assign grouping"),
}

type AssignGrouping struct {
	In  string
	Out GroupSettings
}

// from a list of object names, build a list of group settings
var assignGrouping = testpat.Pattern{
	Name:   "assign grouping",
	Return: "out",
	Labels: []string{"in"},
	Fields: []g.Field{
		{Name: "in", Affinity: affine.Text},
		{Name: "out", Affinity: affine.Record, Type: "group settings"},
	},
	Rules: []rt.Rule{
		{Execute: []rt.Execute{
			&assign.SetValue{
				Target: core.Variable("out", "name"),
				Value:  &assign.FromText{Value: core.Variable("in")}},
			&core.ChooseAction{
				If: &core.Matches{
					Text:  core.Variable("in"),
					Match: "^thing"},
				Does: core.MakeActivity(
					&assign.SetValue{
						Target: core.Variable("out", "label"),
						Value:  &assign.FromText{Value: T("thingies")}},
				),
			},
		}},
	},
}
