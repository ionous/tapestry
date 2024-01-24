package literal_test

import (
	r "reflect"
	"testing"

	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/jsn"
	"git.sr.ht/~ionous/tapestry/lang/encode"
)

// verify that literal commands become plain values.
// for example the object "BoolValue" becomes the bool literal true or false
// note: that slices containing a single value serialize as the value
// this
func TestEncodingDecoding(t *testing.T) {
	testPairs(t, []testPair{{
		&literal.BoolValue{Value: true},
		true,
	}, {
		&literal.NumValue{Value: 5},
		5.0,
	}, {
		&literal.TextValue{Value: "hello"},
		"hello",
	}, {
		&literal.NumValues{Values: []float64{1, 3, 9}},
		[]any{1.0, 3.0, 9.0},
	}, {
		&literal.TextValues{Values: []string{"won", "too", "trees"}},
		[]any{"won", "too", "trees"},
	}, {
		&literal.NumValues{Values: []float64{7}},
		7.0,
	}, {
		&literal.TextValues{Values: []string{"literal"}},
		"literal",
	}})
}

type testPair struct {
	v interface {
		jsn.Marshalee
		literal.LiteralValue
	}
	res any
}

func testPairs(t *testing.T, pairs []testPair) {
	for i, p := range pairs {
		if have, e := marshal(p.v); e != nil {
			t.Log(i, "couldn't encode because", e)
			t.Fail()
		} else if !r.DeepEqual(have, p.res) {
			t.Logf("%d mismatched encode %#v", i, have)
			t.Fail()
		} else {
			// doesnt match the encoding api, but the spirit is the same.
			aff := literal.GetAffinity(p.v)
			if got, e := literal.ReadLiteral(aff, "", p.res); e != nil {
				t.Log(i, "couldn't decode because", e)
				t.Fail()
			} else if !r.DeepEqual(got, p.v) {
				t.Logf("%d mismatched decode %#v", i, have)
				t.Fail()
			}
		}
	}
}

func marshal(v jsn.Marshalee) (ret any, err error) {
	var enc encode.Encoder
	return enc.Customize(literal.CustomEncoder).Encode(v)
}