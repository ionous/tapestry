package test

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/dl/call"
	"git.sr.ht/~ionous/tapestry/dl/list"
	"git.sr.ht/~ionous/tapestry/dl/logic"
	"git.sr.ht/~ionous/tapestry/dl/object"
	"git.sr.ht/~ionous/tapestry/dl/text"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/test/testpat"
)

var runGroupTogther = list.ListMap{
	Target:      object.Variable("settings"),
	List:        &call.FromTextList{Value: object.Variable("objects")},
	PatternName: ("assign grouping"),
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
	Fields: []rt.Field{
		{Name: "in", Affinity: affine.Text},
		{Name: "out", Affinity: affine.Record, Type: "group settings"},
	},
	Rules: []rt.Rule{
		{Exe: []rt.Execute{
			&object.SetValue{
				Target: object.Variable("out", "name"),
				Value:  &call.FromText{Value: object.Variable("in")}},
			&logic.ChooseBranch{
				Condition: &text.Matches{
					Text:  object.Variable("in"),
					Match: "^thing"},
				Exe: []rt.Execute{
					&object.SetValue{
						Target: object.Variable("out", "label"),
						Value:  &call.FromText{Value: T("thingies")}},
				},
			},
		}},
	},
}
