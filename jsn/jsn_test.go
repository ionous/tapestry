package jsn_test

import (
	"encoding/json"
	"hash/fnv"
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
	} else if val := hash(b); val != 0x53398df7 {
		t.Log(val, string(b))
		t.Fatal("mismatched output")
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

func hash(b []byte) uint32 {
	hash := fnv.New32a()
	hash.Write(b)
	return hash.Sum32()
}
