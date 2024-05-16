package core_test

import (
	"encoding/json"
	r "reflect"
	"testing"

	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/dl/core"
	"git.sr.ht/~ionous/tapestry/lang/decode"
	"git.sr.ht/~ionous/tapestry/lang/encode"
	"git.sr.ht/~ionous/tapestry/lang/typeinfo"
	"git.sr.ht/~ionous/tapestry/rt"
	"github.com/kr/pretty"
)

// verify that core variables are written and read as @ strings
func TestCoreEncoding(t *testing.T) {
	testPairs(t, []testPair{{
		&assign.ObjectDot{
			Name: assign.Variable("noun"),
			Dot:  []assign.Dot{&assign.AtField{Field: assign.Variable("trait")}},
		},
		`{"Object:dot:":["@noun",[{"AtField:":"@trait"}]]}`,
	}, {
		&core.AddValue{
			A: assign.Variable("a"),
			B: assign.Object("b", "field"),
		},
		`{"Add:value:":["@a","#b.field"]}`,
	}, {
		// unary
		&core.Softline{},
		`{"Wbr":true}`,
	}, {
		// verify that things that arent variables dont get encoded as variables
		&core.Join{Parts: []rt.TextEval{
			core.T("one"), core.T("two"), core.T("three"),
		}},
		`{"Join parts:":["one","two","three"]}`,
	},
	})
}

type testPair struct {
	// the test serializes this to json
	// to compare against expect
	v typeinfo.Instance
	// the test changes this into json
	// to match against v's json.
	// then unmarshals the json into structs
	// to compares against the original v
	expect string
}

func testPairs(t *testing.T, pairs []testPair) {
	for i, p := range pairs {
		var expect any
		if e := json.Unmarshal([]byte(p.expect), &expect); e != nil {
			t.Fatal(e)
		} else if have, e := marshal(p.v); e != nil {
			t.Errorf("%d couldn't encode because %v", i, e)
		} else if !r.DeepEqual(have, expect) {
			t.Errorf("%d mismatched encode %#v", i, have)
		} else {
			rtype := r.ValueOf(p.v).Elem().Type()
			// println("testing", rtype.String())
			reversed := r.New(rtype).Interface().(typeinfo.Instance)
			if e := unmarshal(reversed, expect); e != nil {
				t.Errorf("%d couldn't decode because %v", i, e)
			} else if !r.DeepEqual(reversed, p.v) {
				t.Errorf("%d mismatched decode", i)
				t.Log("have: ", pretty.Sprint(reversed))
			}
		}
	}
}

func marshal(v typeinfo.Instance) (ret any, err error) {
	var enc encode.Encoder
	return enc.Customize(core.CustomEncoder).Encode(v)
}

func unmarshal(out typeinfo.Instance, plainData any) (err error) {
	var dec decode.Decoder
	return dec.
		Signatures(assign.Z_Types.Signatures, core.Z_Types.Signatures).
		Customize(core.CustomDecoder).
		Decode(out, plainData)
}
