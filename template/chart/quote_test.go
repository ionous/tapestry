package chart

import (
	"testing"
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
	var p QuoteParser
	if e := Parse(&p, str); e != nil {
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
