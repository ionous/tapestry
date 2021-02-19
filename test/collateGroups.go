package test

import (
	"git.sr.ht/~ionous/iffy/affine"
	"git.sr.ht/~ionous/iffy/dl/core"
	"git.sr.ht/~ionous/iffy/dl/list"
	"git.sr.ht/~ionous/iffy/dl/pattern"
	"git.sr.ht/~ionous/iffy/rt"
	g "git.sr.ht/~ionous/iffy/rt/generic"
)

var runCollateGroups = list.Reduce{
	FromList:     V("settings"),
	IntoValue:    "collation",
	UsingPattern: "collate_groups"}

type CollateGroups struct {
	Settings  GroupSettings
	Collation GroupCollation
	Idx       float64
	Groups    []GroupedObjects
	Group     GroupedObjects
	Names     []string
}

var collateGroups = pattern.Pattern{
	Name:   "collate_groups",
	Labels: []string{"settings", "collation"},
	Return: "collation",
	Fields: []g.Field{
		{Name: "settings", Affinity: affine.Record, Type: "group_settings"},
		{Name: "collation", Affinity: affine.Record, Type: "group_collation"},
		{Name: "idx", Affinity: affine.Number},
		{Name: "groups", Affinity: affine.RecordList, Type: "grouped_objects"},
		{Name: "group", Affinity: affine.Record, Type: "grouped_objects"},
		{Name: "names", Affinity: affine.TextList},
	},
	Rules: []*pattern.Rule{
		&pattern.Rule{Execute: core.NewActivity(
			// walk collation.Groups for matching settings
			&core.Assign{
				core.Variable{Str: "groups"},
				// &core.Unpack{V("collation"), "Groups"},
				&core.GetAtField{From: &core.FromVar{N("collation")}, Field: "groups"},
			},
			&list.Each{
				List: V("groups"),
				As:   &list.AsRec{N("el")},
				Do: core.MakeActivity(
					&core.ChooseAction{
						If: &pattern.Determine{
							Pattern: "match_groups",
							Arguments: core.Args(
								V("settings"),
								&core.GetAtField{From: &core.FromVar{N("el")}, Field: "settings"}),
							// &core.Unpack{V("el"), "Settings"}),
						},
						Do: core.MakeActivity(
							&core.Assign{
								Var:  N("idx"),
								From: V("index"),
							},
							// implement a "break" for the each that returns a constant error?
						),
					},
				)}, // end go-each
			&core.ChooseAction{
				If: &core.CompareNum{
					A:  V("idx"),
					Is: &core.EqualTo{},
					B:  &core.Number{0},
				},
				// havent found a matching group?
				// pack the object and its settings into it,
				// push the group into the groups.
				Do: core.MakeActivity(
					&list.PutEdge{
						Into: &list.IntoTxtList{N("names")},
						// From: &core.Unpack{V("settings"), "name"},
						From: &core.GetAtField{From: &core.FromVar{N("settings")}, Field: "name"},
					},
					Put("group", "objects", V("names")),
					Put("group", "settings", V("settings")),
					&list.PutEdge{Into: &list.IntoRecList{N("groups")}, From: V("group")},
				), // end true
				// found a matching group?
				// unpack it, add the object to it, then pack it up again.
				Else: &core.ChooseNothingElse{core.MakeActivity(
					&core.Assign{
						N("group"),
						&core.FromRecord{&list.At{List: V("groups"), Index: V("idx")}}},
					&core.Assign{
						N("names"),
						// &core.Unpack{V("group"), "Objects"},
						&core.GetAtField{From: &core.FromVar{N("group")}, Field: "objects"},
					},
					&list.PutEdge{
						Into: &list.IntoTxtList{N("names")},
						// From: &core.Unpack{V("settings"), "name"},
						From: &core.GetAtField{From: &core.FromVar{N("settings")}, Field: "name"},
					},
					Put("group", "objects", V("names")),
					&list.Set{List: "groups", Index: V("idx"), From: V("group")},
				), // end false
				},
			},
			Put("collation", "groups", V("groups")),
		)},
	},
}

func V(n string) *core.Var {
	return &core.Var{Name: n}
}
func N(n string) core.Variable {
	return core.Variable{Str: n}
}
func Put(rec, field string, from core.Assignment) rt.Execute {
	return &core.PutAtField{Into: &core.IntoVar{Var: N(rec)}, AtField: field, From: from}
}
