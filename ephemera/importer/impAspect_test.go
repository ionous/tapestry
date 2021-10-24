package importer_test

import (
	"testing"

	"git.sr.ht/~ionous/iffy/ephemera/importer"
	"git.sr.ht/~ionous/iffy/ephemera/story"
	"git.sr.ht/~ionous/iffy/jsn"
	"git.sr.ht/~ionous/iffy/jsn/chart"
	"github.com/kr/pretty"
)

// * MarkFlow will have to start passing some memory
// so that we can read ( for instance ) the aspect
// maybe that should be passed to the block functions ( at the very least the block interface )

// then, i think you start chopping down the story importer to highlight what's actually needed.
// and start making the test or something bigger and bigger till its everythings

// mutating slots from one type to another will be a joy im sure.
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
	ts := chart.MakeEncoder(nil)
	if e := ts.Marshal(&asp,
		importer.Map(&ts, story.AspectTraits_Type, importer.KvMap{
			story.AspectTraits_Field_Aspect: func(b jsn.Block, v interface{}) (err error) {
				aspect = *v.(*string)
				return
			},
			story.AspectTraits_Field_TraitPhrase: func(jsn.Block, interface{}) (err error) {
				ts.PushState(importer.EveryValueOf(&ts, story.Trait_Type, func(v interface{}) (err error) {
					trait := *v.(*string)
					out = append(out, aspect+":"+trait)
					return
				}))
				return
			},
		})); e != nil {
		t.Fatal(e)
	} else if diff := pretty.Diff(out, []string{"test:yes", "test:no"}); len(diff) > 0 {
		t.Fatal(diff)
	} else {
		t.Log("okay")
	}
}

// if e := ts.Marshal(&asp, importer.Block(
// 		story.AspectTraits_Type,
// 		func(b jsn.Block) (err error) {
// 			name = asp.Aspect.Str
// 			ts.ChangeState(importer.Eventually(
// 				importer.Value(story.Trait_Type, func(v interface{}) (err error) {
// 					str := *v.(*string)
// 					out = append(out, name+":"+str)
// 					return
// 				})))
// 			return
// 		})); e != nil {
