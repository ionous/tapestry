package assign_test

import (
	"testing"

	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/jsn"
	"git.sr.ht/~ionous/tapestry/jsn/cout"
	"github.com/ionous/errutil"
)

// test that if the user passes various common pattern names through the encoder
// everything works out okay
func TestEncodePattern(t *testing.T) {
	if e := testString("pattern name", "arg name"); e != nil {
		t.Fatal(e)
	} else if e := testString("pattern_name", "arg_name"); e != nil {
		t.Fatal(e)
	} else if e := testString("PatternName", "ArgName"); e != nil {
		t.Fatal(e)
	} else if e := testString("patternName", "argName"); e != nil {
		t.Fatal(e)
	} else if e := testString(" patternName ", " argName  "); e != nil {
		t.Fatal(e)
	}
}
func testString(n, a string) (err error) {
	out := &assign.CallPattern{
		PatternName: n,
		Arguments: []assign.Arg{{
			Name:  a,
			Value: literal.I(5), // the encode gets unhappy without a real value here.1
		}}}

	if have, e := cout.Marshal(out, encoder); e != nil {
		err = e
	} else if have != `{"PatternName argName:":5}` {
		err = errutil.New(have)
	}
	return
}

func encoder(m jsn.Marshaler, flow jsn.FlowBlock) (err error) {
	switch op := flow.GetFlow().(type) {
	case *assign.CallPattern:
		err = assign.EncodePattern(m, op)
	default:
		err = literal.CompactEncoder(m, flow)
	}
	return
}
