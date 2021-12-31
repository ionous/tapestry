package story

import (
	"git.sr.ht/~ionous/tapestry/rt"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
)

// these commands are transformed into other commands at import time.
// fix: replace stubs with custom decoder states?
type StubImporter interface {
	ImportStub(k *Importer) (interface{}, error)
}

var _ StubImporter = (*Comment)(nil)

func (*Comment) Execute(rt.Runtime) error {
	panic("unexpected use of story method")
}

var _ StubImporter = (*CountOf)(nil)

func (*CountOf) GetBool(rt.Runtime) (g.Value, error) {
	panic("unexpected use of story method")
}

var _ StubImporter = (*CycleText)(nil)

func (*CycleText) GetText(rt.Runtime) (g.Value, error) {
	panic("unexpected use of story method")
}

var _ StubImporter = (*Determine)(nil)

func (*Determine) ImportPhrase(k *Importer) (err error) {
	panic("unexpected use of story method")
}

func (*Determine) Execute(rt.Runtime) error {
	panic("unexpected use of story method")
}
func (*Determine) GetBool(rt.Runtime) (g.Value, error) {
	panic("unexpected use of story method")
}
func (*Determine) GetNumber(rt.Runtime) (g.Value, error) {
	panic("unexpected use of story method")
}
func (*Determine) GetText(rt.Runtime) (g.Value, error) {
	panic("unexpected use of story method")
}
func (*Determine) GetRecord(rt.Runtime) (g.Value, error) {
	panic("unexpected use of story method")
}
func (*Determine) GetNumList(rt.Runtime) (g.Value, error) {
	panic("unexpected use of story method")
}
func (*Determine) GetTextList(rt.Runtime) (g.Value, error) {
	panic("unexpected use of story method")
}
func (*Determine) GetRecordList(rt.Runtime) (g.Value, error) {
	panic("unexpected use of story method")
}

var _ StubImporter = (*Make)(nil)

func (*Make) GetRecord(rt.Runtime) (g.Value, error) {
	panic("unexpected use of story method")
}

var _ StubImporter = (*Send)(nil)

func (*Send) GetBool(rt.Runtime) (g.Value, error) {
	panic("unexpected use of story method")
}
func (*Send) Execute(rt.Runtime) error {
	panic("unexpected use of story method")
}

var _ StubImporter = (*ShuffleText)(nil)

func (*ShuffleText) GetText(rt.Runtime) (g.Value, error) {
	panic("unexpected use of story method")
}

var _ StubImporter = (*StoppingText)(nil)

func (*StoppingText) GetText(rt.Runtime) (g.Value, error) {
	panic("unexpected use of story method")
}

var _ StubImporter = (*RenderTemplate)(nil)

func (*RenderTemplate) GetText(rt.Runtime) (g.Value, error) {
	panic("unexpected use of story method")
}
