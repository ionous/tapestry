package core

import (
	"bytes"

	"git.sr.ht/~ionous/iffy/rt"
	g "git.sr.ht/~ionous/iffy/rt/generic"
	"git.sr.ht/~ionous/iffy/rt/print"
)

func (op *Row) GetText(run rt.Runtime) (ret g.Value, err error) {
	// use brackets to establish a span inside the li
	span := print.Brackets("<li>", "</li>")
	return writeSpan(run, span, op, op.Do, span.ChunkOutput())
}

func (op *Rows) GetText(run rt.Runtime) (g.Value, error) {
	var buf bytes.Buffer
	return writeSpan(run, &buf, op, op.Do, print.Tag(&buf, "ul"))
}
