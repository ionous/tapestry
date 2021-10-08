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
	out := dout.NewEncoder()
	src.Marshal(out)
	if d, e := out.Data(); e != nil {
		t.Fatal(e)
	} else if b, e := json.MarshalIndent(d, "", "  "); e != nil {
		t.Fatal(e)
	} else if val := hash(b); val != 0x53398df7 {
		t.Log(string(b))
		t.Fatalf("mismatched output 0x%0x", val)
	} else {
		// var dst story.Story
		// dst.Marshal(din.NewDecoder(b))
		// if diff := pretty.Diff(src, &dst); len(diff) != 0 {
		// 	t.Fatal(diff)
		// }
	}
}

func TestCompact(t *testing.T) {
	src := debug.FactorialStory
	out := cout.NewEncoder()
	src.Marshal(out)
	if d, e := out.Data(); e != nil {
		t.Fatal(e)
	} else if b, e := json.MarshalIndent(d, "", "  "); e != nil {
		t.Fatal(e)
	} else if val := hash(b); val != 0xd86f0fd9 {
		t.Log(string(b))
		t.Fatalf("mismatched output 0x%0x", val)
	} else {

	}
}

func hash(b []byte) uint32 {
	hash := fnv.New32a()
	hash.Write(b)
	return hash.Sum32()
}
