package grok_test

import (
	"reflect"
	"testing"

	"git.sr.ht/~ionous/tapestry/support/grok"
)

// type Word = grok.Word
// type Span = grok.Span
// type MacroType = grok.MacroType
// type Macro = grok.Macro
// type Match = grok.Matched

func TestSep(t *testing.T) {
	cnt := []int{
		grok.Separator(0).Len(),
		grok.CommaSep.Len(),
		grok.AndSep.Len(),
		(grok.CommaSep | grok.AndSep).Len(),
	}
	if !reflect.DeepEqual(cnt, []int{
		0, 1, 1, 2,
	}) {
		t.Fatal(cnt)
	}
}
