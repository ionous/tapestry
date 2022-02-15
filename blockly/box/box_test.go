package box_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"git.sr.ht/~ionous/tapestry/blockly/box"
	"git.sr.ht/~ionous/tapestry/idl"
)

func TestToolbox(t *testing.T) {
	if str, e := box.FromSpecs(idl.Specs); e != nil {
		t.Fatal(e)
	} else if out, e := indent(str); e != nil {
		t.Log(out)
		t.Fatal(e)
	} else {
		t.Log(out)
	}
}

func indent(str string) (ret string, err error) {
	var indent bytes.Buffer
	if e := json.Indent(&indent, []byte(str), "", "  "); e != nil {
		ret = str
		err = e
	} else {
		ret = indent.String()
	}
	return
}
