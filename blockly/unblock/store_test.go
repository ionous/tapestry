package unblock_test

import (
	_ "embed"
	"encoding/json"
	r "reflect"
	"strings"
	"testing"

	"git.sr.ht/~ionous/tapestry"
	"git.sr.ht/~ionous/tapestry/blockly/test"
	"git.sr.ht/~ionous/tapestry/blockly/unblock"
	"git.sr.ht/~ionous/tapestry/jsn"
	"git.sr.ht/~ionous/tapestry/jsn/chart"
	"git.sr.ht/~ionous/tapestry/web/js"
	"github.com/ionous/errutil"
	"github.com/kr/pretty"
)

func TestStoring(t *testing.T) {
	for _, p := range test.Pairs {
		if !strings.HasPrefix(p.Name, "x") {
			t.Run(p.Name, func(t *testing.T) {
				if e := testUnblock(p.Test, p.Json); e != nil {
					t.Fatal(e)
				}
			})
		}
	}
}

// for now just tests that it can load into the in memory structures without error
func xTestUnblockStructs(t *testing.T) {
	var bff unblock.File
	if e := json.Unmarshal(storeTest, &bff); e != nil {
		t.Fatal(e)
	} else {
		// for visual inspection:
		t.Log(pretty.Sprint(bff))
	}
}

//go:embed storeTest.json
var storeTest []byte

// for now just tests that it can load into the in memory structures without error
func TestCountField(t *testing.T) {
	fields := []js.MapItem{
		{Key: "a"},
		{Key: "b0"},
		{Key: "b1"},
		{Key: "b2"},
		{Key: "c"},
	}
	// if i, cnt := (&store.Info{
	// 	Fields: fields,
	// }).CountField("a"); cnt != 0 {
	// 	t.Fatal("a", i, cnt)
	// }
	// if i, cnt := (&store.Info{
	// 	Fields: fields,
	// }).CountField("b"); i != 1 || cnt != 3 {
	// 	t.Fatal("b", i, cnt)
	// }
	if i, cnt := (&unblock.Info{
		Fields: fields[1:2],
	}).CountField("b"); i != 0 || cnt != 1 {
		t.Fatal("b", i, cnt)
	}
}

func testUnblock(expect jsn.Marshalee, msg string) (err error) {
	var top unblock.Info
	if e := json.Unmarshal([]byte(msg), &top); e != nil {
		err = e
	} else {
		dst := r.New(r.TypeOf(expect).Elem()).Interface().(jsn.Marshalee)
		dec := chart.MakeDecoder()
		if e := dec.Marshal(dst, unblock.NewTopBlock(&dec, tapestry.Registry(), &top)); e != nil {
			err = e
		} else if diff := pretty.Diff(expect, dst); len(diff) > 0 {
			pretty.Println(dst)
			err = errutil.New(e, "mismatched", diff)
		}
	}
	return
}
