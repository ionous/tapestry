package test

import (
	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/jsn"
)

var Pairs = []struct {
	Name string
	Test jsn.Marshalee
	Json string
}{{
	// a primitive list is a list of dummy inputs
	// noting that blockly ignores dummies when saving --
	// so they get saved in the "fields" section
	"List",
	&literal.TextValues{
		Values: []string{"a", "b", "c"},
	}, `{
  "id": "test-1",
  "type": "text_values",
  "extraState": {
    "VALUES": 3
  },
  "fields": {
    "VALUES0": "a",
    "VALUES1": "b",
    "VALUES2": "c"
  }
}`,
}}
