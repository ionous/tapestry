package charmed

import (
	"strings"
	"testing"

	"git.sr.ht/~ionous/tapestry/support/charm"
	"github.com/ionous/errutil"
)

func TestSpaces(t *testing.T) {
	// index of the fail point, or -1 if success is expected
	count := func(failPoint int, str string, style charm.State) (err error) {
		if e := charm.Parse(style, str); e == nil && failPoint != -1 {
			err = errutil.New("unexpected success")
		} else if n, ok := e.(charm.EndpointError); !ok {
			err = e
		} else if at := n.End(); at != failPoint {
			// 0 means okay, -1 incomplete, >0 the one-index of the failure point.
			err = errutil.New(str, "len:", at)
		}
		return
	}
	if e := count(0, "a", RequiredSpaces); e != nil {
		t.Log("ok: e")
	} else {
		t.Fatal("expected failure")
	}
	if e := count(0, "a", OptionalSpaces); e != nil {
		t.Fatal(e)
	}
	if e := count(-1, strings.Repeat(" ", 5), OptionalSpaces); e != nil {
		t.Fatal(e)
	}
	if e := count(3, strings.Repeat(" ", 3)+"x", OptionalSpaces); e != nil {
		t.Fatal(e)
	}
}
