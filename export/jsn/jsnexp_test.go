package jsn_test

import (
	"encoding/json"
	"testing"

	"git.sr.ht/~ionous/iffy/ephemera/debug"
	"git.sr.ht/~ionous/iffy/export/jsn"
)

func TestDetails(t *testing.T) {
	src := debug.FactorialStory
	dm := jsn.NewDetailedMarshaler()
	src.Marshal(dm)
	if b, e := json.MarshalIndent(dm.Data(), "", "  "); e != nil {
		t.Fatal(e)
	} else {
		t.Log(string(b))
	}
}
