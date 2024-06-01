package shortcut_test

import (
	"encoding/json"
	r "reflect"
	"testing"

	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/dl/math"
	"git.sr.ht/~ionous/tapestry/dl/object"
	"git.sr.ht/~ionous/tapestry/dl/printer"
	"git.sr.ht/~ionous/tapestry/dl/text"
	"git.sr.ht/~ionous/tapestry/lang/decode"
	"git.sr.ht/~ionous/tapestry/lang/encode"
	"git.sr.ht/~ionous/tapestry/lang/shortcut"
	"git.sr.ht/~ionous/tapestry/lang/typeinfo"
	"git.sr.ht/~ionous/tapestry/rt"
	"github.com/kr/pretty"
)

// verify that core variables are written and read as @ strings
func TestCoreEncoding(t *testing.T) {
	testPairs(t, []testPair{{
		&object.ObjectDot{
			Name: object.Variable("noun"),
			Dot:  []object.Dot{&object.AtField{Field: object.Variable("trait")}},
		},
		`{"Object:dot:":["@noun",[{"AtField:":"@trait"}]]}`,
	}, {
		&math.AddValue{
			A: object.Variable("a"),
			B: object.Object("b", "field"),
		},
		`{"Add:value:":["@a","#b.field"]}`,
	}, {
		// unary
		&printer.Softline{},
		`{"Wbr":true}`,
	}, {
		// verify that things that arent variables dont get encoded as variables
		&text.Join{Parts: []rt.TextEval{
			literal.T("one"), literal.T("two"), literal.T("three"),
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
	return enc.Customize(shortcut.Encoder).Encode(v)
}

func unmarshal(out typeinfo.Instance, plainData any) (err error) {
	var dec decode.Decoder
	return dec.
		Signatures(
			assign.Z_Types.Signatures,
			math.Z_Types.Signatures,
			object.Z_Types.Signatures,
			printer.Z_Types.Signatures,
			text.Z_Types.Signatures,
		).
		Customize(shortcut.Decoder).
		Decode(out, plainData)
}
