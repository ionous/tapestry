package story_test

import (
	"database/sql"
	"log"
	"testing"

	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/dl/core"
	"git.sr.ht/~ionous/tapestry/dl/list"
	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/dl/story"
	"git.sr.ht/~ionous/tapestry/qna"
	"git.sr.ht/~ionous/tapestry/qna/decode"
	"git.sr.ht/~ionous/tapestry/qna/qdb"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
	"git.sr.ht/~ionous/tapestry/test/testweave"
	"git.sr.ht/~ionous/tapestry/weave/assert"
)

// generate ephemera for macros
func TestMacros(t *testing.T) {
	// ugh. this setup.
	dt := testweave.NewWeaverOptions(t.Name(), func(db *sql.DB) rt.Runtime {
		qx, e := qdb.NewQueryx(db)
		if e != nil {
			panic(e)
		}
		return qna.NewRuntimeOptions(
			log.Writer(),
			qx,
			decode.NewDecoder(story.AllSignatures),
			qna.NewOptions(),
		)
	}, false)
	defer dt.Close()
	cat := dt.Catalog()
	//
	if e := cat.AssertDomainStart("tapestry", nil); e != nil {
		t.Fatal(e)
	} else if e := addDefaultKinds(cat); e != nil {
		t.Fatal(e)
	} else if e := story.ImportStory(cat, t.Name(), macroStory); e != nil {
		t.Fatal(e)
	} else if e := cat.AssertDomainEnd(); e != nil {
		t.Fatal(e)
	} else if e := cat.AssembleCatalog(); e != nil {
		t.Fatal(e)
	} else {
		// expect := []eph.Ephemera{
		// 	&eph.Values{
		// 		Noun:  "Hershel",
		// 		Field: "proper_named",
		// 		Value: literal.B(true),
		// 	},
		// 	&eph.Nouns{
		// 		Noun: "Hershel",
		// 		Kind: "an actor",
		// 	},
		// 	&eph.Nouns{
		// 		Noun: "scissors",
		// 		Kind: "things",
		// 	},
		// 	&eph.Relatives{
		// 		Rel:       "whereabouts",
		// 		Noun:      "Hershel",
		// 		OtherNoun: "scissors",
		// 	},
	}
	// if diff := pretty.Diff(els, expect); len(diff) != 0 {
	// 	t.Log(pretty.Sprint(els))
	// 	t.Fatal(diff)
	// }
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
		// &story.KindOfRelation{
		// },
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
