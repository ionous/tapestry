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

// Say some bit of text.
// FIX: should take a speaker, and we should have a default speaker
type Say struct {
	Text rt.TextEval `if:"selector"`
}

// Buffer collects text said by other statements and returns them as a string.
// Unlike Span, it does not add or alter spaces between writes.
type Buffer struct {
	Do Activity
}

// Span collects text printed during a block and writes the text with spaces.
type Span struct {
	Do Activity
}

// Bracket sandwiches text printed during a block and puts them inside parenthesis ().
type Bracket struct {
	Do Activity
}

// Slash separates text printed during a block with left-leaning slashes.
type Slash struct {
	Do Activity
}

// Commas writes words separated with commas, ending with an "and".
type Commas struct {
	Do Activity
}

// Compose defines a spec for the composer editor.
func (*Say) Compose() composer.Spec {
	return composer.Spec{
		Name:   "say_text",
		Group:  "printing",
		Desc:   "Say: print some bit of text to the player.",
		Fluent: &composer.Fluid{Role: composer.Command},
	}
}

// Execute writes text to the runtime's current writer.
func (op *Say) Execute(run rt.Runtime) (err error) {
	if e := safe.WriteText(run, op.Text); e != nil {
		err = cmdError(op, e)
	}
	return
}

// Compose defines a spec for the composer editor.
func (*Buffer) Compose() composer.Spec {
	return composer.Spec{
		Name:  "buffer_text",
		Group: "printing",
	}
}

func (op *Buffer) GetText(run rt.Runtime) (g.Value, error) {
	var buf bytes.Buffer
	return writeSpan(run, &buf, op, op.Do, &buf)
}

// Compose defines a spec for the composer editor.
func (*Span) Compose() composer.Spec {
	return composer.Spec{
		Name:  "span_text",
		Group: "printing",
		Desc:  "Span Text: Writes text with spaces between words.",
	}
}

func (op *Span) GetText(run rt.Runtime) (g.Value, error) {
	span := print.NewSpanner() // separate writes with spaces
	return writeSpan(run, span, op, op.Do, span.ChunkOutput())
}

// Compose defines a spec for the composer editor.
func (*Bracket) Compose() composer.Spec {
	return composer.Spec{
		Name:  "bracket_text",
		Group: "printing",
		Desc:  "Bracket text: Sandwiches text printed during a block and puts them inside parenthesis '()'.",
	}
}

func (op *Bracket) GetText(run rt.Runtime) (g.Value, error) {
	span := print.Parens()
	return writeSpan(run, span, op, op.Do, span.ChunkOutput())
}

// Compose defines a spec for the composer editor.
func (*Slash) Compose() composer.Spec {
	return composer.Spec{
		Name:  "slash_text",
		Group: "printing",
		Desc:  "Slash text: Separates words with left-leaning slashes '/'.",
	}
}

func (op *Slash) GetText(run rt.Runtime) (g.Value, error) {
	span := print.NewSpanner() // separate punctuation with spaces
	return writeSpan(run, span, op, op.Do, print.Slash(span.ChunkOutput()))
}

// Compose defines a spec for the composer editor.
func (*Commas) Compose() composer.Spec {
	return composer.Spec{
		Name:  "comma_text",
		Group: "printing",
		Desc:  "Comma text: Separates words with commas, and 'and'.",
	}
}

func (op *Commas) GetText(run rt.Runtime) (g.Value, error) {
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
				ret = g.StringOf(safe.HackTillTemplatesCanEvaluatePatternTypes)
			}
			safe.HackTillTemplatesCanEvaluatePatternTypes = ""
		}
	}
	return
}
