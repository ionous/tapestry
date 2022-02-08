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
	return writeSpan(run, span, op, op.Does, span.ChunkOutput())
}

func (op *Rows) GetText(run rt.Runtime) (g.Value, error) {
	var buf bytes.Buffer
	return writeSpan(run, &buf, op, op.Does, print.Tag(&buf, "ul"))
}
