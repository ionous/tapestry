package js

import (
	"encoding/json"
	"testing"
)

func TestMapSlice(t *testing.T) {
	var els MapSlice
	const data = `{"0" : null, 
	"1": -0.1E+4, 
	"boop": "s", "3": [ null, 0, "string", [], {
}], 
	"last": {"0": null, "1": 0, "2": "s", "3": [], "4": {}}
}`
	if e := json.Unmarshal([]byte(data), &els); e != nil {
		t.Fatal(e)
	} else {
		var keys = [...]string{
			"0", "1", "boop", "3", "last",
		}
		var vals = [...]string{
			`null`,
			`-0.1E+4`,
			`"s"`,
			"[ null, 0, \"string\", [], {\n}]",
			`{"0": null, "1": 0, "2": "s", "3": [], "4": {}}`,
		}
		var x [1]struct{} // static assert that the lengths are the same
		var _ = x[len(keys)-len(vals)]
		//
		if want, got := len(keys), len(els); want != got {
			t.Fatal("expected length", want, "got", got)
		}
		for i, el := range els {
			if want, got := keys[i], el.Key; want != got {
				t.Fatal("expected key", want, "got", got)
			}
			if want, got := vals[i], string(el.Msg); want != got {
				t.Fatal("expected value", want, "got", got)
			}
		}
	}
}
