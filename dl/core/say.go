package core

import (
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/print"
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"git.sr.ht/~ionous/tapestry/rt/writer"
)

// Execute writes text to the runtime's current writer.
func (op *PrintText) Execute(run rt.Runtime) (err error) {
	if e := safe.WriteText(run, op.Text); e != nil {
		err = cmdError(op, e)
	}
	return
}

func (op *PrintWords) Execute(run rt.Runtime) error {
	return safe.WriteText(run, op)
}

func (op *PrintWords) GetText(run rt.Runtime) (ret rt.Value, err error) {
	var span print.Spanner
	if v, e := writeSpan(run, &span, op.Exe, span.ChunkOutput()); e != nil {
		err = cmdError(op, e)
	} else {
		ret = v
	}
	return
}

func (op *PrintParens) Execute(run rt.Runtime) error {
	return safe.WriteText(run, op)
}

func (op *PrintParens) GetText(run rt.Runtime) (ret rt.Value, err error) {
	span := print.Parens()
	if v, e := writeSpan(run, &span, op.Exe, span.ChunkOutput()); e != nil {
		err = cmdError(op, e)
	} else {
		ret = v
	}
	return
}

func (op *PrintCommas) Execute(run rt.Runtime) error {
	return safe.WriteText(run, op)
}

func (op *PrintCommas) GetText(run rt.Runtime) (ret rt.Value, err error) {
	var span print.Spanner // separate punctuation with spaces
	if v, e := writeSpan(run, &span, op.Exe, print.AndSeparator(span.ChunkOutput())); e != nil {
		err = cmdError(op, e)
	} else {
		ret = v
	}
	return
}

type stringer interface{ String() string }

func writeSpan(run rt.Runtime, buf stringer, act []rt.Execute, out writer.OutputCloser) (ret rt.Value, err error) {
	if len(act) == 0 {
		ret = rt.Empty
	} else {
		prev := run.SetWriter(out)
		if e := safe.RunAll(run, act); e != nil {
			err = e
		} else if e := out.Close(); e != nil {
			err = e
		} else if str := buf.String(); len(str) > 0 {
			ret = rt.StringOf(str)
		} else {
			ret = safe.GetTemplateText()
		}
		run.SetWriter(prev)
	}
	return
}
