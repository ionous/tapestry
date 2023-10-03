package box_test

import (
	"encoding/json"
	"testing"

	"git.sr.ht/~ionous/tapestry/blockly/box"
	"git.sr.ht/~ionous/tapestry/idl"
)

func TestToolbox(t *testing.T) {
	if str, e := box.FromSpecs(idl.Specs); e != nil {
		t.Fatal(e)
	} else if !json.Valid([]byte(str)) {
		t.Fatal(str)
	} else {
		t.Log(str)
	}
}
