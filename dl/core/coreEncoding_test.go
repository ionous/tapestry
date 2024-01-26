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
		&assign.ObjectRef{
			Name:  core.Variable("noun"),
			Field: core.Variable("trait"),
		},
		`{"Object:field:":["@noun","@trait"]}`,
	}, {
		// fix:  should have path syntax ( re: expressions )  ex. @pawn.trait
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
	v      typeinfo.Inspector
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
			// println("testing", rtype.String())
			reversed := r.New(rtype).Interface().(typeinfo.Inspector)
			if e := unmarshal(reversed, expect); e != nil {
				t.Logf("%d couldn't decode because %v", i, e)
				t.Fail()
			} else if !r.DeepEqual(reversed, p.v) {
				t.Logf("%d mismatched decode", i)
				t.Log("want: ", pretty.Sprint(p.v))
				t.Log("have: ", pretty.Sprint(reversed))
				t.Fail()
			}
		}
	}
}

func marshal(v typeinfo.Inspector) (ret any, err error) {
	var enc encode.Encoder
	return enc.Customize(core.CustomEncoder).Encode(v)
}

func unmarshal(out typeinfo.Inspector, plainData any) (err error) {
	var dec decode.Decoder
	return dec.
		Signatures(assign.Z_Types.Signatures, core.Z_Types.Signatures).
		Customize(core.CustomDecoder).
		Decode(out, plainData)
}

func readPlainData(str string) (ret any) {
	if e := json.Unmarshal([]byte(str), &ret); e != nil {
		panic(e)
	}
	return
}
