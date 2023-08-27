package core

import (
	"bytes"

	"git.sr.ht/~ionous/tapestry/rt"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"git.sr.ht/~ionous/tapestry/rt/print"
)

func (op *Row) GetText(run rt.Runtime) (ret g.Value, err error) {
	// use brackets to establish a span inside the li
	span := print.Brackets("<li>", "</li>")
	if v, e := writeSpan(run, &span, op.Exe, span.ChunkOutput()); e != nil {
		err = cmdError(op, e)
	} else {
		ret = v
	}
	return
}

func (op *Rows) GetText(run rt.Runtime) (ret g.Value, err error) {
	var buf bytes.Buffer
	if v, e := writeSpan(run, &buf, op.Exe, print.Tag(&buf, "ul")); e != nil {
		err = cmdError(op, e)
	} else {
		ret = v
	}
	return
}
