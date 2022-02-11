package block_test

import (
  "bytes"
  "encoding/json"
  "strconv"
  "testing"

  "git.sr.ht/~ionous/tapestry/blockly/block"
  "git.sr.ht/~ionous/tapestry/blockly/test"
  "git.sr.ht/~ionous/tapestry/dl/literal"
  "git.sr.ht/~ionous/tapestry/dl/testdl"
  "git.sr.ht/~ionous/tapestry/jsn"
  "git.sr.ht/~ionous/tapestry/jsn/chart"
  "git.sr.ht/~ionous/tapestry/web/js"
  "github.com/ionous/errutil"
  "github.com/kr/pretty"
)

// test a flow within a flow
func TestEmbeds(t *testing.T) {
  if e := testBlocks(&testdl.TestEmbed{
    TestFlow: testdl.TestFlow{},
  }, `{
  "id": "test-1",
  "type": "test_embed",
  "extraState": {
    "TEST_FLOW": 1
  },
  "inputs": {
    "TEST_FLOW": {
      "block": {
        "id": "test-2",
        "type": "test_flow",
        "extraState": {}
      }
    }
  }
}`); e != nil {
    t.Fatal(e)
  }
}

// test a swap member of the flow
func TestSwap(t *testing.T) {
  if e := testBlocks(&testdl.TestFlow{
    Swap: testdl.TestSwap{
      Choice: testdl.TestSwap_C_Opt,
      Value: &testdl.TestTxt{
        Str: "something",
      },
    },
  }, `{
  "id": "test-1",
  "type": "test_flow",
  "extraState": {
    "SWAP": 1
  },
  "fields": {
    "SWAP": "$C"
  },
  "inputs": {
    "SWAP": {
      "block": {
        "id": "test-2",
        "type": "test_txt",
        "fields": {
          "TEST_TXT": "something"
        }
      }
    }
  }
}`); e != nil {
    t.Fatal(e)
  }
}

// test a slot member of the flow
func TestSlot(t *testing.T) {
  if e := testBlocks(&literal.FieldValue{
    Field: "test",
    Value: &literal.NumValue{
      Num: 5,
    }}, `{
  "id": "test-1",
  "type": "field_value",
  "extraState": {
    "FIELD": 1,
    "VALUE": 1
  },
  "fields": {
    "FIELD": "test"
  },
  "inputs": {
    "VALUE": {
      "block": {
        "id": "test-2",
        "type": "num_value",
        "extraState": {
          "NUM": 1
        },
        "fields": {
          "NUM": 5
        }
      }
    }
  }
}`); e != nil {
    t.Fatal(e)
  }
}

func TestPairs(t *testing.T) {
  for _, p := range test.Pairs {
    t.Run(p.Name, func(t *testing.T) {
      if e := testBlocks(p.Test, p.Json); e != nil {
        t.Fatal(e)
      }
    })
  }
}

// repeats of a specific flow
func TestSlice(t *testing.T) {
  if e := testBlocks(&literal.FieldValues{
    Contains: []literal.FieldValue{{
      Field: "first",
      Value: &literal.NumValue{
        Num: 5,
      }}, {
      Field: "second",
      Value: &literal.TextValue{
        Text: "five",
      }},
    },
  }, `{
  "id": "test-1",
  "type": "field_values",
  "extraState": {
    "CONTAINS": 2
  },
  "inputs": {
    "CONTAINS0": {
      "block": {
        "id": "test-2",
        "type": "field_value",
        "extraState": {
          "FIELD": 1,
          "VALUE": 1
        },
        "fields": {
          "FIELD": "first"
        },
        "inputs": {
          "VALUE": {
            "block": {
              "id": "test-3",
              "type": "num_value",
              "extraState": {
                "NUM": 1
              },
              "fields": {
                "NUM": 5
              }
            }
          }
        }
      }
    },
    "CONTAINS1": {
      "block": {
        "id": "test-4",
        "type": "field_value",
        "extraState": {
          "FIELD": 1,
          "VALUE": 1
        },
        "fields": {
          "FIELD": "second"
        },
        "inputs": {
          "VALUE": {
            "block": {
              "id": "test-5",
              "type": "text_value",
              "extraState": {
                "TEXT": 1
              },
              "fields": {
                "TEXT": "five"
              }
            }
          }
        }
      }
    }
  }
}`); e != nil {
    t.Fatal(e)
  }
}

// repeats of a non-stacking slot.
func TestSeries(t *testing.T) {
  if e := testBlocks(&testdl.TestFlow{
    Slots: []testdl.TestSlot{
      &testdl.TestFlow{},
      &testdl.TestFlow{},
    }}, `{
  "id": "test-1",
  "type": "test_flow",
  "extraState": {
    "SLOTS": 2
  },
  "inputs": {
    "SLOTS0": {
      "block": {
        "id": "test-2",
        "type": "test_flow",
        "extraState": {}
      }
    },
    "SLOTS1": {
      "block": {
        "id": "test-3",
        "type": "test_flow",
        "extraState": {}
      }
    }
  }
}`); e != nil {
    t.Fatal(e)
  }
}

// repeats of an empty series
func TestEmptySeries(t *testing.T) {
  if e := testBlocks(&testdl.TestFlow{
    Slots: []testdl.TestSlot{},
  }, `{
  "id": "test-1",
  "type": "test_flow",
  "extraState": {}
}`); e != nil {
    t.Fatal(e)
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
