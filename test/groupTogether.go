package test

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/dl/core"
	"git.sr.ht/~ionous/tapestry/dl/list"
	"git.sr.ht/~ionous/tapestry/rt"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"git.sr.ht/~ionous/tapestry/test/testpat"
)

var runGroupTogther = list.ListMap{
	Target:       core.Variable("settings"),
	List:         core.AssignFromTextList(GetVariable("objects")),
	UsingPattern: W("assign_grouping"),
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
			&core.SetValue{
				Target: core.Variable("out", "name"),
				Value:  core.AssignFromText(GetVariable("in"))},
			&core.ChooseAction{
				If: &core.Matches{
					Text:    GetVariable("in"),
					Pattern: "^thing"},
				Does: core.MakeActivity(
					&core.SetValue{
						Target: core.Variable("out", "label"),
						Value:  core.AssignFromText(T("thingies"))},
				),
			},
		}},
	},
}
