package unblock_test

import (
	_ "embed"
	"encoding/json"
	r "reflect"
	"strings"
	"testing"

	"git.sr.ht/~ionous/tapestry/blockly/test"
	"git.sr.ht/~ionous/tapestry/blockly/unblock"
	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/dl/story"
	"git.sr.ht/~ionous/tapestry/dl/testdl"
	"git.sr.ht/~ionous/tapestry/lang/typeinfo"
	"git.sr.ht/~ionous/tapestry/web/js"
	"github.com/ionous/errutil"
	"github.com/kr/pretty"
)

func TestStoring(t *testing.T) {
	reg := unblock.MakeBlockCreator([]*typeinfo.TypeSet{
		&literal.Z_Types,
		&testdl.Z_Types,
		&story.Z_Types, // for Comment; fix: make a test_empty?
	})

	for _, p := range test.Pairs {
		if strings.HasPrefix(p.Name, "x") {
			t.Log("skipping", p.Name)
		} else {
			t.Run(p.Name, func(t *testing.T) {
				if e := testUnblock(reg, p.Test, p.Json); e != nil {
					t.Fatal(e)
				}
			})
		}
	}
}

// for now just tests that it can load into the in memory structures without error
func TestCountField(t *testing.T) {
	fields := []js.MapItem{
		{Key: "a"},
		{Key: "b0"},
		{Key: "b1"},
		{Key: "b2"},
		{Key: "c"},
	}
	// if i, cnt := (&store.BlockInfo{
	// 	Fields: fields,
	// }).CountField("a"); cnt != 0 {
	// 	t.Fatal("a", i, cnt)
	// }
	// if i, cnt := (&store.BlockInfo{
	// 	Fields: fields,
	// }).CountField("b"); i != 1 || cnt != 3 {
	// 	t.Fatal("b", i, cnt)
	// }
	if i, cnt := (&unblock.BlockInfo{
		Fields: fields[1:2],
	}).CountFields("b"); i != 0 || cnt != 1 {
		t.Fatal("b", i, cnt)
	}
}

func testUnblock(reg unblock.Creator, expect typeinfo.Instance, msg string) (err error) {
	var top unblock.BlockInfo
	if e := json.Unmarshal([]byte(msg), &top); e != nil {
		err = e
	} else {
		ptr := r.New(r.TypeOf(expect).Elem()).Interface().(typeinfo.Instance)
		if e := unblock.DecodeBlock(ptr, reg, &top); e != nil {
			err = e
		} else {
			if diff := pretty.Diff(expect, ptr); len(diff) > 0 {
				pretty.Println(ptr)
				err = errutil.New(e, "mismatched", diff)
			}
		}
	}
	return
}
