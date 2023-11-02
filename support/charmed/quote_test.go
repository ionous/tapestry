package charmed

import (
	"testing"

	"git.sr.ht/~ionous/tapestry/support/charm"
	"github.com/ionous/errutil"
)

func TestQuotes(t *testing.T) {
	if x, e := testQ(t, "'singles'", "singles"); e != nil {
		t.Fatal(x, e)
	}
	if x, e := testQ(t, `"doubles"`, "doubles"); e != nil {
		t.Fatal(x, e)
	}
	if x, e := testQ(t, "'escape\"'", "escape\""); e != nil {
		t.Fatal(x, e)
	}
	if x, e := testQ(t, `"\\"`, "\\"); e != nil {
		t.Fatal(x, e)
	}
	if x, e := testQ(t, string([]rune{'"', '\\', 'a', '"'}), "\a"); e != nil {
		t.Fatal(x, e)
	}
	if _, e := testQ(t, string([]rune{'"', '\\', 'g', '"'}),
		ignoreResult); e == nil {
		t.Fatal(e)
	}
}

func testQ(t *testing.T, str, want string) (ret interface{}, err error) {
	t.Log("test:", str)
	var p quoteParser
	if e := charm.ParseEof(&p, str); e != nil {
		err = e
	} else if got, e := p.GetString(); e != nil {
		err = e
	} else if want != ignoreResult {
		if got != want {
			err = mismatched(want, got)
		} else {
			t.Log("ok", got)
		}
	}
	return str, err
}

func mismatched(want, got string) error {
	return errutil.Fmt("want(%d): %s; != got(%d): %s.", len(want), want, len(got), got)
}

// for testing errors when we want to fail before the match is tested.
const ignoreResult = "~~IGNORE~~"

type quoteParser struct {
	QuoteParser
}

// NewRune starts with the leading quote mark; it finishes just after the matching quote mark.
func (p *QuoteParser) NewRune(r rune) (ret charm.State) {
	if r == '\'' || r == '"' {
		ret = p.ScanQuote(r)
	}
	return
}
