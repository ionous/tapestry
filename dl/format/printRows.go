package format

import (
	"bytes"

	"git.sr.ht/~ionous/tapestry/dl/cmd"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/print"
	"git.sr.ht/~ionous/tapestry/rt/safe"
)

func (op *PrintRow) Execute(run rt.Runtime) (err error) {
	return safe.WriteText(run, op)
}

func (op *PrintRow) GetText(run rt.Runtime) (ret rt.Value, err error) {
	// use brackets to establish a span inside the li
	span := print.Brackets("<li>", "</li>")
	if v, e := writeSpan(run, &span, op.Exe, span.ChunkOutput()); e != nil {
		err = cmd.Error(op, e)
	} else {
		ret = v
	}
	return
}

func (op *PrintRows) Execute(run rt.Runtime) (err error) {
	return safe.WriteText(run, op)
}

func (op *PrintRows) GetText(run rt.Runtime) (ret rt.Value, err error) {
	var buf bytes.Buffer
	if v, e := writeSpan(run, &buf, op.Exe, print.Tag(&buf, "ul")); e != nil {
		err = cmd.Error(op, e)
	} else {
		ret = v
	}
	return
}
