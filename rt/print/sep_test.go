package print

import (
	"io"
	"testing"

	"git.sr.ht/~ionous/tapestry/rt/writer"
)

func TestPrintSep(t *testing.T) {
	if s, e := write(AndSeparator, "pizza"); e != nil {
		t.Fatal(e)
	} else if expect := "pizza"; s != expect {
		t.Fatalf("mismatched want (%v), have (%v)", expect, s)
	}
	if s, e := write(AndSeparator, "apple", "hedgehog", "washington", "mushroom"); e != nil {
		t.Fatal(e)
	} else if expect := "apple, hedgehog, washington, and mushroom"; s != expect {
		t.Fatalf("mismatched want (%v), have (%v)", expect, s)
	}
	if s, e := write(AndSeparator, "apple", "hedgehog"); e != nil {
		t.Fatal(e)
	} else if expect := "apple and hedgehog"; s != expect {
		t.Fatalf("serial comma only after two items; want (%v), have (%v)", expect, s)
	}
	//
	if s, e := write(OrSeparator, "pistachio"); e != nil {
		t.Fatal(e)
	} else if expect := "pistachio"; s != expect {
		t.Fatalf("mismatched want (%v), have (%v)", expect, s)
	}
	if s, e := write(OrSeparator, "apple", "hedgehog", "washington", "mushroom"); e != nil {
		t.Fatal(e)
	} else if expect := "apple, hedgehog, washington, or mushroom"; s != expect {
		t.Fatalf("mismatched want (%v), have (%v)", expect, s)
	}
	if s, e := write(OrSeparator, "washington", "mushroom"); e != nil {
		t.Fatal(e)
	} else if expect := "washington or mushroom"; s != expect {
		t.Fatalf("serial comma only after two items, mismatched want (%v), have (%v)", expect, s)
	}
}

func write(sep func(w writer.Output) writer.OutputCloser, names ...string) (ret string, err error) {
	span := NewSpanner()
	w := sep(span.ChunkOutput())
	for _, n := range names {
		if _, e := io.WriteString(w, n); e != nil {
			err = e
			break
		}
	}
	if err == nil {
		// normally PopWriter would call close, but we arent using the runtime here.
		if e := w.Close(); e != nil {
			err = e
		} else {
			ret = span.String()
		}
	}
	return
}
