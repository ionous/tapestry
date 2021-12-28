package core

import (
	"bytes"

	"git.sr.ht/~ionous/iffy/dl/composer"
	"git.sr.ht/~ionous/iffy/rt"
	g "git.sr.ht/~ionous/iffy/rt/generic"
	"git.sr.ht/~ionous/iffy/rt/print"
	"git.sr.ht/~ionous/iffy/rt/safe"
	"git.sr.ht/~ionous/iffy/rt/writer"
	"github.com/ionous/errutil"
)

// Execute writes text to the runtime's current writer.
func (op *Say) Execute(run rt.Runtime) (err error) {
	if e := safe.WriteText(run, op.Text); e != nil {
		err = cmdError(op, e)
	}
	return
}

func (op *BufferText) GetText(run rt.Runtime) (g.Value, error) {
	var buf bytes.Buffer
	return writeSpan(run, &buf, op, op.Do, &buf)
}

func (op *SpanText) GetText(run rt.Runtime) (g.Value, error) {
	span := print.NewSpanner() // separate writes with spaces
	return writeSpan(run, span, op, op.Do, span.ChunkOutput())
}

func (op *BracketText) GetText(run rt.Runtime) (g.Value, error) {
	span := print.Parens()
	return writeSpan(run, span, op, op.Do, span.ChunkOutput())
}

func (op *SlashText) GetText(run rt.Runtime) (g.Value, error) {
	span := print.NewSpanner() // separate punctuation with spaces
	return writeSpan(run, span, op, op.Do, print.Slash(span.ChunkOutput()))
}

func (op *CommaText) GetText(run rt.Runtime) (g.Value, error) {
	span := print.NewSpanner() // separate punctuation with spaces
	return writeSpan(run, span, op, op.Do, print.AndSeparator(span.ChunkOutput()))
}

type stringer interface{ String() string }

// s - access to what was written
// op - for reporting errors
// act - activity that presumably generates some output
// w - output target with any needed filters, etc.
// returns the output of "s" as a value
func writeSpan(run rt.Runtime, s stringer, op composer.Composer, act Activity, w writer.Output) (ret g.Value, err error) {
	if !act.Empty() {
		was := run.SetWriter(w)
		ex := act.Execute(run)
		run.SetWriter(was)
		if e := errutil.Append(ex, writer.Close(w)); e != nil {
			err = cmdError(op, e)
		} else {
			if res := s.String(); len(res) > 0 {
				ret = g.StringOf(res)
			} else {
				ret = safe.HackTillTemplatesCanEvaluatePatternTypes
			}
			safe.HackTillTemplatesCanEvaluatePatternTypes = nil
		}
	}
	return
}
