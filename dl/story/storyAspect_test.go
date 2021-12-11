package story_test

import (
	"testing"

	"git.sr.ht/~ionous/iffy/dl/story"
	"git.sr.ht/~ionous/iffy/jsn"
	"git.sr.ht/~ionous/iffy/jsn/chart"
	"github.com/kr/pretty"
)

// FIX: mutating slots from one type to another will be a joy im sure.
func TestEventually(t *testing.T) {
	asp := story.AspectTraits{
		Aspect: story.Aspect{Str: "test"},
		TraitPhrase: story.TraitPhrase{
			Trait: []story.Trait{{
				Str: "yes",
			}, {
				Str: "no",
			}},
		},
	}
	var aspect string
	var out []string
	ts := chart.MakeEncoder()
	if e := ts.Marshal(&asp,
		story.Map(&ts, story.BlockMap{
			story.AspectTraits_Type: story.KeyMap{
				story.AspectTraits_Field_Aspect: func(b jsn.Block, v interface{}) (err error) {
					aspect = *v.(*string)
					return
				},
				story.AspectTraits_Field_TraitPhrase: func(jsn.Block, interface{}) (err error) {
					ts.PushState(story.EveryValueOf(&ts, story.Trait_Type, func(v interface{}) (err error) {
						trait := *v.(*string)
						out = append(out, aspect+":"+trait)
						return
					}))
					return
				},
			}})); e != nil {
		t.Fatal(e)
	} else if diff := pretty.Diff(out, []string{"test:yes", "test:no"}); len(diff) > 0 {
		t.Fatal(diff)
	} else {
		t.Log("okay")
	}
}

func TestEndBlock(t *testing.T) {
	asp := story.Certainties{
		PluralKinds: story.PluralKinds{Str: "test"},
	}
	found := false
	ts := chart.MakeEncoder()
	if e := ts.Marshal(&asp,
		story.Map(&ts, story.BlockMap{
			story.Certainties_Type: story.KeyMap{
				story.BlockEnd: func(b jsn.Block, v interface{}) (err error) {
					cs := b.(jsn.FlowBlock).GetFlow().(*story.Certainties) // ick
					found = cs.PluralKinds.Str == "test"
					return
				},
			}})); e != nil {
		t.Fatal(e)
	} else if !found {
		t.Fatal("end not found")
	}
}
