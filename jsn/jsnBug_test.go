package jsn_test

import (
	"testing"

	"git.sr.ht/~ionous/tapestry/dl/story"
	"git.sr.ht/~ionous/tapestry/rt"
	"github.com/kr/pretty"
)

// previously being unable to parse the contents of a pattern call would fail to throw an error
func TestBug_PatternOfPattern(t *testing.T) {
	str := `{"PrintVantage where:":{"ParentOf obj:":"@actor"}}`
	var dst rt.Execute
	if e := story.DecodeJson(rt.Execute_Slot{&dst}, []byte(str), nil); e == nil {
		pretty.Println(dst)
		t.Fatal("expected failure")
	}
}
