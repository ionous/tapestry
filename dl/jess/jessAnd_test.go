package jess_test

import (
	"reflect"
	"testing"

	"git.sr.ht/~ionous/tapestry/dl/jess"
)

// verify the returned length of various separator values
func TestCommaLength(t *testing.T) {
	cnt := []int{
		jess.Separator(0).Len(),
		jess.CommaSep.Len(),
		jess.AndSep.Len(),
		(jess.CommaSep | jess.AndSep).Len(),
	}
	if !reflect.DeepEqual(cnt, []int{
		0, 1, 1, 2,
	}) {
		t.Fatal(cnt)
	}
}
