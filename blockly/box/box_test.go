package box_test

import (
	"encoding/json"
	"testing"

	"git.sr.ht/~ionous/tapestry/blockly/box"
	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/dl/call"
	"git.sr.ht/~ionous/tapestry/dl/debug"
	"git.sr.ht/~ionous/tapestry/dl/format"
	"git.sr.ht/~ionous/tapestry/dl/frame"
	"git.sr.ht/~ionous/tapestry/dl/game"
	"git.sr.ht/~ionous/tapestry/dl/grammar"
	"git.sr.ht/~ionous/tapestry/dl/jess"
	"git.sr.ht/~ionous/tapestry/dl/list"
	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/dl/logic"
	"git.sr.ht/~ionous/tapestry/dl/math"
	"git.sr.ht/~ionous/tapestry/dl/object"
	"git.sr.ht/~ionous/tapestry/dl/prim"
	"git.sr.ht/~ionous/tapestry/dl/rel"
	"git.sr.ht/~ionous/tapestry/dl/render"
	"git.sr.ht/~ionous/tapestry/dl/rtti"
	"git.sr.ht/~ionous/tapestry/dl/story"
	"git.sr.ht/~ionous/tapestry/dl/text"
	"git.sr.ht/~ionous/tapestry/lang/typeinfo"
	"git.sr.ht/~ionous/tapestry/support/files"
)

func TestToolbox(t *testing.T) {
	if str, e := box.FromTypes(blocks); e != nil {
		t.Fatal(e)
	} else if !json.Valid([]byte(str)) {
		t.Fatal(str)
	} else {
		str := files.Indent(str)
		t.Log(str)
	}
}

var blocks = []*typeinfo.TypeSet{
	&assign.Z_Types,
	&call.Z_Types,
	&debug.Z_Types,
	&frame.Z_Types,
	&game.Z_Types,
	&grammar.Z_Types,
	&jess.Z_Types,
	&list.Z_Types,
	&literal.Z_Types,
	&logic.Z_Types,
	&math.Z_Types,
	&object.Z_Types,
	// &play.Z_Types,
	&prim.Z_Types,
	&format.Z_Types,
	&rel.Z_Types,
	&render.Z_Types,
	&game.Z_Types,
	&rtti.Z_Types,
	&story.Z_Types,
	// &testdl.Z_Types,
	&text.Z_Types,
}
