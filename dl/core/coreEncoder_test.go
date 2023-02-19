package core_test

import (
	"testing"

	"git.sr.ht/~ionous/tapestry/dl/core"
	"git.sr.ht/~ionous/tapestry/jsn/cout"
)

func TestCompactEncoder(t *testing.T) {
	out := &core.CallPattern{
		Pattern: core.P("factorial"),
		Arguments: []core.Arg{{
			Name: "num",
			/// in the old bits this will generate an error if the assignment is emptyl
			Value: core.AssignFromNumber(core.I(1)),
		}}}

	if have, e := cout.Marshal(out, core.CompactEncoder); e != nil {
		t.Fatal(e)
	} else if have != `{"Factorial num:":{"FromNumber:":1}}` {
		t.Fatal(have)
	}
}
