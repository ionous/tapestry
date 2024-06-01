package text

import (
	"testing"

	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/test/testutil"
)

func TestMatches(t *testing.T) {
	var run testutil.PanicRuntime
	// test a valid regexp
	// loop to verify(ish) the cache
	m := &Matches{Text: literal.T("gophergopher"), Match: "(gopher){2}"}
	for i := 0; i < 2; i++ {
		if ok, e := m.GetBool(&run); e != nil {
			t.Fatal(e)
		} else if !ok.Bool() {
			t.Fatal("mismatch")
		}
	}
	// test that a bad expression will fail
	// loop to verify(ish) the cache
	fail := &Matches{Match: "("}
	for i := 0; i < 2; i++ {
		if _, e := fail.GetBool(&run); e == nil {
			t.Fatal("expected error")
		} else if _, e := fail.GetBool(&run); e == nil {
			t.Fatal("expected error")
		}
	}
}
