package jsn_test

import (
	"encoding/json"
	"testing"

	"git.sr.ht/~ionous/iffy/ephemera/debug"
	"git.sr.ht/~ionous/iffy/jsn/cout"
	"git.sr.ht/~ionous/iffy/jsn/dout"
)

func TestDetails(t *testing.T) {
	src := debug.FactorialStory
	m := dout.NewDetailedMarshaler()
	src.Marshal(m)
	if d, e := m.Data(); e != nil {
		t.Fatal(e)
	} else if b, e := json.MarshalIndent(d, "", "  "); e != nil {
		t.Fatal(e)
	} else {
		t.Log(string(b))
	}
}

func TestCompact(t *testing.T) {
	src := debug.FactorialStory
	m := cout.NewCompactMarshaler()
	src.Marshal(m)
	if d, e := m.Data(); e != nil {
		t.Fatal(e)
	} else if b, e := json.MarshalIndent(d, "", "  "); e != nil {
		t.Fatal(e)
	} else {
		t.Log(string(b))
	}
}
