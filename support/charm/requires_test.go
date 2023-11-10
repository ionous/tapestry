package charm

import (
	"strings"
	"testing"

	"github.com/ionous/errutil"
)

func TestRequires(t *testing.T) {
	isSpace := func(r rune) bool { return r == ' ' }

	// index of the fail point, or -1 if success is expected
	count := func(failPoint int, str string, style State) (err error) {
		if e := ParseEof(style, str); e == nil && failPoint != -1 {
			err = errutil.New("unexpected success")
		} else if n, ok := e.(EndpointError); !ok {
			err = e
		} else if at := n.End(); at != failPoint {
			// 0 means okay, -1 incomplete, >0 the one-index of the failure point.
			err = errutil.New(str, "len:", at)
		}
		return
	}
	if e := count(0, "a", AtleastOne(isSpace)); e != nil {
		t.Fatal(e)
	}
	if e := count(0, "a", Optional(isSpace)); e != nil {
		t.Fatal(e)
	}
	if e := count(-1, strings.Repeat(" ", 5), Optional(isSpace)); e != nil {
		t.Fatal(e)
	}
	if e := count(3, strings.Repeat(" ", 3)+"x", Optional(isSpace)); e != nil {
		t.Fatal(e)
	}
}
