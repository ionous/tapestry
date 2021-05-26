package core

import (
	"io"

	"git.sr.ht/~ionous/iffy/rt"
)

func (op *Blankline) Execute(run rt.Runtime) (_ error) {
	io.WriteString(run.Writer(), "<p>")
	return
}

func (op *Newline) Execute(run rt.Runtime) (_ error) {
	io.WriteString(run.Writer(), "<br>")
	return
}

func (op *Softline) Execute(run rt.Runtime) (_ error) {
	io.WriteString(run.Writer(), "<wbr>")
	return
}
