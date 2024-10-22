package rt_test

import (
	"testing"

	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/rt"
)

func TestSplices(t *testing.T) {
	zeroSplice := func(src rt.Value) {
		var vs rt.Value
		if e := src.Splice(0, 0, nil, &vs); e != nil {
			t.Fatal("empty splice should be legal")
		} else if vs == nil {
			t.Fatal("empty splice should return value")
		} else if a := vs.Affinity(); a != affine.TextList {
			t.Fatal("empty splice should return a text list not", a)
		} else if cnt := vs.Len(); cnt != 0 {
			t.Fatal("empty splice should return an empty list", cnt)
		}
	}
	zeroSplice(rt.StringsOf(nil))
	zeroSplice(rt.StringsOf([]string{"a"}))
	zeroSplice(rt.StringsOf([]string{"a", "b", "c"}))
}
