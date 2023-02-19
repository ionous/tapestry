package core

import (
	"bytes"

	"git.sr.ht/~ionous/tapestry/rt"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"git.sr.ht/~ionous/tapestry/rt/print"
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"git.sr.ht/~ionous/tapestry/rt/writer"
	"github.com/ionous/errutil"
)

// Execute writes text to the runtime's current writer.
func (op *Say) Execute(run rt.Runtime) (err error) {
	if e := safe.WriteText(run, op.Text); e != nil {
		err = cmdError(op, e)
	}
	return
}

func (op *BufferText) GetText(run rt.Runtime) (ret g.Value, err error) {
	var buf bytes.Buffer
	if v, e := writeSpan(run, &buf, op.Does, &buf); e != nil {
		err = cmdError(op, e)
	} else {
		ret = v
	}
	return
}

func (op *SpanText) GetText(run rt.Runtime) (ret g.Value, err error) {
	span := print.NewSpanner() // separate writes with spaces
	if v, e := writeSpan(run, span, op.Does, span.ChunkOutput()); e != nil {
		err = cmdError(op, e)
	} else {
		ret = v
	}
	return
}

func (op *BracketText) GetText(run rt.Runtime) (ret g.Value, err error) {
	span := print.Parens()
	if v, e := writeSpan(run, span, op.Does, span.ChunkOutput()); e != nil {
		err = cmdError(op, e)
	} else {
		ret = v
	}
	return
}

func (op *SlashText) GetText(run rt.Runtime) (ret g.Value, err error) {
	span := print.NewSpanner() // separate punctuation with spaces
	if v, e := writeSpan(run, span, op.Does, print.Slash(span.ChunkOutput())); e != nil {
		err = cmdError(op, e)
	} else {
		ret = v
	}
	return
}

func (op *CommaText) GetText(run rt.Runtime) (ret g.Value, err error) {
	span := print.NewSpanner() // separate punctuation with spaces
	if v, e := writeSpan(run, span, op.Does, print.AndSeparator(span.ChunkOutput())); e != nil {
		err = cmdError(op, e)
	} else {
		ret = v
	}
	return
}

type stringer interface{ String() string }

// s - access to what was written
// op - for reporting errors
// act - activity that presumably generates some output
// w - output target with any needed filters, etc.
// returns the output of "s" as a value
func writeSpan(run rt.Runtime, s stringer, act []rt.Execute, w writer.Output) (ret g.Value, err error) {
	if len(act) == 0 {
		ret = g.Empty
	} else {
		ret, err = WriteSpan(run, w, s, func() error {
			return safe.RunAll(run, act)
		})
	}
	return
}

func WriteSpan(run rt.Runtime, w writer.Output, s stringer, cb func() error) (ret g.Value, err error) {
	was := run.SetWriter(w)
	ex := cb()
	run.SetWriter(was)
	if e := errutil.Append(ex, writer.Close(w)); e != nil {
		err = e
	} else {
		if res := s.String(); len(res) > 0 {
			ret = g.StringOf(res)
		} else if hack := safe.HackTillTemplatesCanEvaluatePatternTypes; hack != nil {
			// we didn't accumulate any text during execution
			// but perhaps we ran a pattern that returned text.
			// to get rid of this, we'd examine (at runtime or compile time) the futures calls
			// and switch on execute patterns vs text patterns
			// an example is { .Lantern } which says the name
			// vs. { pluralize: .Lantern } which returns the pluralized name.
			ret = hack
			safe.HackTillTemplatesCanEvaluatePatternTypes = nil
		} else {
			ret = g.Empty // if the res was empty, it might have intentionally been empty
		}
	}
	return
}
