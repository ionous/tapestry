package story_test

import (
	"testing"

	"git.sr.ht/~ionous/tapestry"
	"git.sr.ht/~ionous/tapestry/dl/core"
	"git.sr.ht/~ionous/tapestry/dl/eph"
	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/dl/story"
	"git.sr.ht/~ionous/tapestry/imp"
	"git.sr.ht/~ionous/tapestry/jsn"
	"git.sr.ht/~ionous/tapestry/support/asm"
	"git.sr.ht/~ionous/tapestry/test/debug"

	"git.sr.ht/~ionous/tapestry/jsn/cout"
)

func TestImportStory(t *testing.T) {
	var curr story.Story
	if e := story.Decode(&curr, debug.Blob, tapestry.AllSignatures); e != nil {
		t.Fatal(e)
	} else {
		var els []eph.Ephemera
		k := imp.NewImporter(collectEphemera(&els), storyMarshaller)
		if e := asm.ImportStory(k, t.Name(), &curr); e != nil {
			t.Fatal("import", e)
		} else {
			t.Log("ok")
		}
	}
}

func B(b bool) *literal.BoolValue   { return &literal.BoolValue{Value: b} }
func I(n int) *literal.NumValue     { return &literal.NumValue{Value: float64(n)} }
func F(n float64) *literal.NumValue { return &literal.NumValue{Value: n} }
func T(s string) *literal.TextValue { return &literal.TextValue{Value: s} }

func P(p string) core.PatternName  { return core.PatternName{Str: p} }
func N(v string) core.VariableName { return core.VariableName{Str: v} }
func V(i string) *core.GetVar      { return &core.GetVar{Name: N(i)} }
func W(v string) string            { return v }

func storyMarshaller(m jsn.Marshalee) (string, error) {
	return cout.Marshal(m, story.CompactEncoder)
}

func collectEphemera(sink *[]eph.Ephemera) imp.WriterFun {
	return func(el eph.Ephemera) {
		*sink = append(*sink, el)
	}
}
