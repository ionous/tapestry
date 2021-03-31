package core

import (
	"io"

	"git.sr.ht/~ionous/iffy/dl/composer"
	"git.sr.ht/~ionous/iffy/rt"
)

type Paragraph struct{}
type Newline struct{}
type Softline struct{}

// Compose returns a spec for use by the composer editor.
func (*Paragraph) Compose() composer.Spec {
	return composer.Spec{
		Name:   "p",
		Group:  "printing",
		Desc:   "Paragraph: add a single blank line following some text.",
		Fluent: &composer.Fluid{Name: "p", Role: composer.Command},
	}
}

func (*Newline) Compose() composer.Spec {
	return composer.Spec{
		Name:   "br",
		Group:  "printing",
		Desc:   "Newline: start a new line.",
		Fluent: &composer.Fluid{Name: "br", Role: composer.Command},
	}
}

func (*Softline) Compose() composer.Spec {
	return composer.Spec{
		Name:   "wbr",
		Group:  "printing",
		Desc:   "Softline: start a new line ( if not already at a new line. )",
		Fluent: &composer.Fluid{Name: "wbr", Role: composer.Command},
	}
}

func (op *Paragraph) Execute(run rt.Runtime) (_ error) {
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
