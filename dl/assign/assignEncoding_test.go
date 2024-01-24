package assign_test

import (
	"reflect"
	"testing"

	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/lang/encode"
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
	var enc encode.Encoder
	out := &assign.CallPattern{
		PatternName: n,
		Arguments: []assign.Arg{{
			Name: a,
			Value: &assign.FromNumber{
				Value: literal.I(5), // the encode gets unhappy without a real value here.
			},
		}}}
	// calls EncodePattern indirectly
	if got, e := enc.Customize(assign.CustomEncoder).Encode(out); e != nil {
		err = e
	} else if !reflect.DeepEqual(got, wantPattern) {
		err = errutil.New("mismatch", got)
	}
	return
}

var wantPattern = map[string]any{
	"PatternName argName:": map[string]any{
		"FromNumber:": 5.0,
	},
}