package block_test

import (
  "bytes"
  "encoding/json"
  "strconv"
  "testing"

  "git.sr.ht/~ionous/tapestry/blockly/block"
  "git.sr.ht/~ionous/tapestry/blockly/test"
  "git.sr.ht/~ionous/tapestry/jsn"
  "git.sr.ht/~ionous/tapestry/jsn/chart"
  "git.sr.ht/~ionous/tapestry/web/js"
  "github.com/ionous/errutil"
  "github.com/kr/pretty"
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
  if e := enc.Marshal(src, block.NewTopBlock(&enc, &out)); e != nil {
    err = errutil.New(e, "failed marshal")
  } else if str, e := indent(out.String()); e != nil {
    err = errutil.New(e, "invalid json", str)
  } else if diff := pretty.Diff(str, expect); len(diff) > 0 {
    println(str)
    err = errutil.New(e, "mismatched", diff, str)
  }
  return
}

func indent(str string) (ret string, err error) {
  var indent bytes.Buffer
  if e := json.Indent(&indent, []byte(str), "", "  "); e != nil {
    err = e
    ret = str
  } else {
    ret = indent.String()
  }
  return
}
