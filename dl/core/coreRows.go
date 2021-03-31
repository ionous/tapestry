package core

import (
	"bytes"

	"git.sr.ht/~ionous/iffy/dl/composer"
	"git.sr.ht/~ionous/iffy/rt"
	g "git.sr.ht/~ionous/iffy/rt/generic"
	"git.sr.ht/~ionous/iffy/rt/print"
)

type Row struct {
	Do Activity
}

type Rows struct {
	Do Activity
}

func (*Row) Compose() composer.Spec {
	return composer.Spec{
		Name:   "row",
		Group:  "printing",
		Desc:   "Row: a single line as part of a group of lines.",
		Fluent: &composer.Fluid{Name: "row", Role: composer.Command},
	}
}

// Compose returns a spec for use by the composer editor.
func (*Rows) Compose() composer.Spec {
	return composer.Spec{
		Name:   "rows",
		Group:  "printing",
		Desc:   "Rows: group text into successive lines.",
		Fluent: &composer.Fluid{Name: "rows", Role: composer.Command},
	}
}

func (op *Row) GetText(run rt.Runtime) (ret g.Value, err error) {
	// use brackets to establish a span inside the li
	span := print.Brackets("<li>", "</li>")
	return writeSpan(run, span, op, op.Do, span.ChunkOutput())
}

func (op *Rows) GetText(run rt.Runtime) (g.Value, error) {
	var buf bytes.Buffer
	return writeSpan(run, &buf, op, op.Do, print.Tag(&buf, "ul"))
}
