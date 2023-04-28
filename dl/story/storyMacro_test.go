package story_test

import (
	"testing"

	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/dl/core"
	"git.sr.ht/~ionous/tapestry/dl/list"
	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/dl/story"
	"git.sr.ht/~ionous/tapestry/imp/assert"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
)

// generate ephemera for macros
func xTestMacros(t *testing.T) {
	// errutil.Panic = true
	// //
	// var cat eph.Catalog
	// k := imp.NewImporter(cat)
	// if e := k.BeginDomain("tapestry", nil); e != nil {
	// 	t.Fatal(e)
	// } else if e := addDefaultKinds(k); e != nil {
	// 	t.Fatal(e)
	// } else if e := story.ImportStory(k, t.Name(), macroStory); e != nil {
	// 	t.Fatal(e)
	// } else if e := cat.AssembleCatalog(); e != nil {
	// 	t.Fatal(e)
	// } else {
	// 	// expect := []eph.Ephemera{
	// 	// 	&eph.EphValues{
	// 	// 		Noun:  "Hershel",
	// 	// 		Field: "proper_named",
	// 	// 		Value: literal.B(true),
	// 	// 	},
	// 	// 	&eph.EphNouns{
	// 	// 		Noun: "Hershel",
	// 	// 		Kind: "an actor",
	// 	// 	},
	// 	// 	&eph.EphNouns{
	// 	// 		Noun: "scissors",
	// 	// 		Kind: "things",
	// 	// 	},
	// 	// 	&eph.EphRelatives{
	// 	// 		Rel:       "whereabouts",
	// 	// 		Noun:      "Hershel",
	// 	// 		OtherNoun: "scissors",
	// 	// 	},
	// 	// }
	// 	// if diff := pretty.Diff(els, expect); len(diff) != 0 {
	// 	// 	t.Log(pretty.Sprint(els))
	// 	// 	t.Fatal(diff)
	// 	// }
	// }
}

func addDefaultKinds(n assert.Assertions) (err error) {
	for _, k := range kindsOf.DefaultKinds {
		if e := n.AssertAncestor(k.String(), k.Parent().String()); e != nil {
			err = e
			break
		}
	}
	return
}

var macroStory = &story.StoryFile{
	StoryStatements: []story.StoryStatement{
		&story.DefineMacro{
			MacroName: core.T("carrier"),
			Params: []story.FieldDefinition{
				// todo: optional labels for fields so these can be "carrier", "nouns" internally
				// ( and maybe that would allow the anonymous first parameter )
				&story.TextField{
					Name: "actor",
				},
				&story.TextListField{
					Name: "carries",
				},
			},
			Result: &story.NothingField{},
			Locals: nil,
			MacroStatements: []rt.Execute{
				// assert:
				// 1. the actor is an actor
				&story.DefineNouns{
					// fix: autoconversions of text to text list?
					Nouns: &list.MakeTextList{
						// tbd: allow variables to have determiners? (and strip them off during import or weave)
						Values: []rt.TextEval{core.Variable("actor")},
					},
					Kind: literal.T("an actor"),
				},
				// 2. each thing is a thing ( the relation is any object so... )
				&story.DefineNouns{
					Nouns: core.Variable("carries"),
					Kind:  literal.T("things"),
				},
				// 3. use whereabouts
				&story.DefineRelatives{
					Nouns:    core.Variable("carries"),
					Relation: literal.T("whereabouts"),
					OtherNouns: &list.MakeTextList{
						Values: []rt.TextEval{core.Variable("actor")},
					},
				},
			},
		},
		&story.DefineScene{
			Scene: literal.T("testing"),
			With: []story.StoryStatement{
				&story.CallMacro{
					MacroName: "carrier",
					Arguments: []assign.Arg{{
						Name:  "actor",
						Value: &assign.FromText{Value: literal.T("Hershel")},
					}, {
						Name:  "carries",
						Value: &assign.FromTextList{Value: literal.Ts("the scissors")},
					}}},
			},
		},
	},
}
