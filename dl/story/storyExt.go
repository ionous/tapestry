package story

import (
	"git.sr.ht/~ionous/iffy/rt"
	g "git.sr.ht/~ionous/iffy/rt/generic"
)

type StubImporter interface {
	ImportStub(k *Importer) (interface{}, error)
}

var _ StubImporter = (*Comment)(nil)

func (*Comment) Execute(rt.Runtime) error {
	panic("unexpected use of story method in runtime")
}

var _ StubImporter = (*CountOf)(nil)

func (*CountOf) GetBool(rt.Runtime) (g.Value, error) {
	panic("unexpected use of story method in runtime")
}

var _ StubImporter = (*CycleText)(nil)

func (*CycleText) GetText(rt.Runtime) (g.Value, error) {
	panic("unexpected use of story method in runtime")
}

var _ StubImporter = (*Determine)(nil)

func (*Determine) Execute(rt.Runtime) error {
	panic("unexpected use of story method in runtime")
}
func (*Determine) GetBool(rt.Runtime) (g.Value, error) {
	panic("unexpected use of story method in runtime")
}
func (*Determine) GetNumber(rt.Runtime) (g.Value, error) {
	panic("unexpected use of story method in runtime")
}
func (*Determine) GetText(rt.Runtime) (g.Value, error) {
	panic("unexpected use of story method in runtime")
}
func (*Determine) GetRecord(rt.Runtime) (g.Value, error) {
	panic("unexpected use of story method in runtime")
}
func (*Determine) GetNumList(rt.Runtime) (g.Value, error) {
	panic("unexpected use of story method in runtime")
}
func (*Determine) GetTextList(rt.Runtime) (g.Value, error) {
	panic("unexpected use of story method in runtime")
}
func (*Determine) GetRecordList(rt.Runtime) (g.Value, error) {
	panic("unexpected use of story method in runtime")
}

var _ StubImporter = (*Make)(nil)

func (*Make) GetRecord(rt.Runtime) (g.Value, error) {
	panic("unexpected use of story method in runtime")
}

var _ StubImporter = (*Send)(nil)

func (*Send) GetBool(rt.Runtime) (g.Value, error) {
	panic("unexpected use of story method in runtime")
}
func (*Send) Execute(rt.Runtime) error {
	panic("unexpected use of story method in runtime")
}

var _ StubImporter = (*ShuffleText)(nil)

func (*ShuffleText) GetText(rt.Runtime) (g.Value, error) {
	panic("unexpected use of story method in runtime")
}

var _ StubImporter = (*StoppingText)(nil)

func (*StoppingText) GetText(rt.Runtime) (g.Value, error) {
	panic("unexpected use of story method in runtime")
}

var _ StubImporter = (*RenderTemplate)(nil)

func (*RenderTemplate) GetText(rt.Runtime) (g.Value, error) {
	panic("unexpected use of story method in runtime")
}
