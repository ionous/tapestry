package core

import (
	"io"

	"git.sr.ht/~ionous/tapestry/rt"
)

func (op *ParagraphBreak) Execute(run rt.Runtime) (_ error) {
	io.WriteString(run.Writer(), "<p>")
	return
}

func (op *LineBreak) Execute(run rt.Runtime) (_ error) {
	io.WriteString(run.Writer(), "<br>")
	return
}

func (op *Softline) Execute(run rt.Runtime) (_ error) {
	io.WriteString(run.Writer(), "<wbr>")
	return
}
