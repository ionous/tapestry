package core

import (
	"bytes"

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

// collect all text written during this function and return it as a value
func (op *BufferText) GetText(run rt.Runtime) (ret rt.Value, err error) {
	var buf bytes.Buffer
	if v, e := writeSpan(run, &buf, op.Exe, &writer.ChunkWriter{Writer: &buf}); e != nil {
		err = cmdError(op, e)
	} else {
		ret = v
	}
	return
}

func (op *SpanText) GetText(run rt.Runtime) (ret rt.Value, err error) {
	var span print.Spanner
	if v, e := writeSpan(run, &span, op.Exe, span.ChunkOutput()); e != nil {
		err = cmdError(op, e)
	} else {
		ret = v
	}
	return
}

func (op *BracketText) GetText(run rt.Runtime) (ret rt.Value, err error) {
	span := print.Parens()
	if v, e := writeSpan(run, &span, op.Exe, span.ChunkOutput()); e != nil {
		err = cmdError(op, e)
	} else {
		ret = v
	}
	return
}

func (op *SlashText) GetText(run rt.Runtime) (ret rt.Value, err error) {
	var span print.Spanner // separate punctuation with spaces
	if v, e := writeSpan(run, &span, op.Exe, print.Slash(span.ChunkOutput())); e != nil {
		err = cmdError(op, e)
	} else {
		ret = v
	}
	return
}

func (op *CommaText) GetText(run rt.Runtime) (ret rt.Value, err error) {
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
