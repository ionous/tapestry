package core_test

import (
	"encoding/json"
	r "reflect"
	"testing"

	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/dl/core"
	"git.sr.ht/~ionous/tapestry/jsn"
	"git.sr.ht/~ionous/tapestry/lang/decode"
	"git.sr.ht/~ionous/tapestry/lang/encode"
)

// verify that core variables are written and read as @ strings
// fix: these should have path syntax ( re: expressions )
// ex. @pawn.trait
func TestEncodingDecoding(t *testing.T) {
	testPairs(t, []testPair{{
		&assign.ObjectRef{
			Name:  core.Variable("noun"),
			Field: core.Variable("trait"),
		},
		`{"Object:field:":["@noun","@trait"]}`,
	}, {
		core.Variable("pawn", "trait"),
		`{"Variable:dot:":["pawn",[{"AtField:":"trait"}]]}`,
	}, {
		&core.AddValue{
			A: core.Variable("a"),
			B: core.Variable("b"),
		},
		`{"Add:value:":["@a","@b"]}`,
	}, {
		// unary
		&core.Softline{},
		`{"Wbr":true}`,
	},
	})
}

type testPair struct {
	v interface {
		jsn.Marshalee
	}
	expect string
}

func testPairs(t *testing.T, pairs []testPair) {
	for i, p := range pairs {
		var expect any
		if e := json.Unmarshal([]byte(p.expect), &expect); e != nil {
			t.Fatal(e)
		} else if have, e := marshal(p.v); e != nil {
			t.Logf("%d couldn't encode because %v", i, e)
			t.Fail()
		} else if !r.DeepEqual(have, expect) {
			t.Logf("%d mismatched encode %#v", i, have)
			t.Fail()
		} else {
			rtype := r.ValueOf(p.v).Elem().Type()
			println("newing", rtype.String())
			reversed := r.New(rtype).Interface().(jsn.Marshalee)
			if e := unmarshal(reversed, expect); e != nil {
				t.Logf("%d couldn't decode because %v", i, e)
				t.Fail()
			} else if !r.DeepEqual(reversed, p.v) {
				t.Logf("%d mismatched decode %#v", i, have)
				t.Fail()
			}
		}
	}
}

func marshal(v jsn.Marshalee) (ret any, err error) {
	var enc encode.Encoder
	return enc.Customize(core.CustomEncoder).Encode(v)
}
func unmarshal(out jsn.Marshalee, plainData any) (err error) {
	var dec decode.Decoder
	return dec.
		Signatures(assign.Signatures, core.Signatures).
		Customize(core.CustomDecoder).
		Decode(out, plainData)
}

func readPlainData(str string) (ret any) {
	if e := json.Unmarshal([]byte(str), &ret); e != nil {
		panic(e)
	}
	return
}
