package block_test

import (
	"strconv"
	"testing"

	"git.sr.ht/~ionous/tapestry/blockly/block"
	"git.sr.ht/~ionous/tapestry/blockly/test"
	"git.sr.ht/~ionous/tapestry/lang/typeinfo"
	"git.sr.ht/~ionous/tapestry/support/files"
	"git.sr.ht/~ionous/tapestry/web/js"
	"github.com/ionous/errutil"
)

func TestPairs(t *testing.T) {
	for _, p := range test.Pairs {
		t.Run(p.Name, func(t *testing.T) {
			if e := testBlocks(p.Test, p.Json); e != nil {
				t.Fatal(e)
			}
		})
	}
}

func testBlocks(src typeinfo.Instance, expect string) (err error) {
	var id int
	block.NewId = func() string {
		id++
		return "test-" + strconv.Itoa(id)
	}
	var out js.Builder
	if e := block.Build(&out, src, false); e != nil {
		err = errutil.New(e, "failed marshal")
	} else if str := files.Indent(out.String()); str != expect {
		err = errutil.New("mismatched", str)
	}
	return
}
