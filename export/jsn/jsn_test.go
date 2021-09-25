package jsn_test

import (
	"encoding/json"
	"testing"

	"git.sr.ht/~ionous/iffy/ephemera/debug"
	"git.sr.ht/~ionous/iffy/export/jsn/compact"
	"git.sr.ht/~ionous/iffy/export/jsn/detailed"
)

func TestDetails(t *testing.T) {
	src := debug.FactorialStory
	m := detailed.NewDetailedMarshaler()
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
	m := compact.NewCompactMarshaler()
	src.Marshal(m)
	if d, e := m.Data(); e != nil {
		t.Fatal(e)
	} else if b, e := json.MarshalIndent(d, "", "  "); e != nil {
		t.Fatal(e)
	} else {
		t.Log(string(b))
	}
}
