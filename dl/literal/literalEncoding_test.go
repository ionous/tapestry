package literal_test

import (
	r "reflect"
	"testing"

	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/lang/encode"
	"git.sr.ht/~ionous/tapestry/lang/typeinfo"
)

// verify that literal commands become plain values.
// for example the object "BoolValue" becomes the bool literal true or false
// note: that slices containing a single value serialize as the value
// this
func TestLiteralEncoding(t *testing.T) {
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
		&literal.NumList{Values: []float64{1, 3, 9}},
		[]any{1.0, 3.0, 9.0},
	}, {
		&literal.TextList{Values: []string{"won", "too", "trees"}},
		[]any{"won", "too", "trees"},
	}, {
		&literal.NumList{Values: []float64{7}},
		7.0,
	}, {
		&literal.TextList{Values: []string{"literal"}},
		"literal",
	}})
}

type testPair struct {
	v interface {
		typeinfo.Instance
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

func marshal(v typeinfo.Instance) (ret any, err error) {
	var enc encode.Encoder
	return enc.Customize(literal.CustomEncoder).Encode(v)
}
