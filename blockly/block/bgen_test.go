package block_test

import (
	"strconv"
	"testing"

	"git.sr.ht/~ionous/tapestry/blockly/block"
	"git.sr.ht/~ionous/tapestry/blockly/test"
	"git.sr.ht/~ionous/tapestry/dl/spec/rs"
	"git.sr.ht/~ionous/tapestry/idl"
	"git.sr.ht/~ionous/tapestry/jsn"
	"git.sr.ht/~ionous/tapestry/jsn/chart"
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

func testBlocks(src jsn.Marshalee, expect string) (err error) {
	var id int
	block.NewId = func() string {
		id++
		return "test-" + strconv.Itoa(id)
	}
	var out js.Builder
	enc := chart.MakeEncoder()
	if ts, e := rs.FromSpecs(idl.Specs); e != nil {
		err = e
	} else if e := enc.Marshal(src, block.NewTopBlock(&enc, &ts, &out, false)); e != nil {
		err = errutil.New(e, "failed marshal")
	} else if str := jsn.Indent(out.String()); str != expect {
		err = errutil.New(e, "mismatched", str)
	}
	return
}
