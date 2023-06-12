package story

import (
	"git.sr.ht/~ionous/tapestry/rt"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"git.sr.ht/~ionous/tapestry/weave"
)

// these commands are transformed into other commands at import time.
// fix: replace stubs with custom decoder states?

// PreImport happens at the opening of a json block and it can transform the value into something completely new.
type PreImport interface {
	PreImport(*weave.Catalog) (interface{}, error)
}

var _ PreImport = (*CountOf)(nil)

func (*CountOf) GetBool(rt.Runtime) (g.Value, error) {
	panic("unexpected use of story method")
}

var _ PreImport = (*CycleText)(nil)

func (*CycleText) GetText(rt.Runtime) (g.Value, error) {
	panic("unexpected use of story method")
}

var _ PreImport = (*ShuffleText)(nil)

func (*ShuffleText) GetText(rt.Runtime) (g.Value, error) {
	panic("unexpected use of story method")
}

var _ PreImport = (*StoppingText)(nil)

func (*StoppingText) GetText(rt.Runtime) (g.Value, error) {
	panic("unexpected use of story method")
}

// the import step transforms say template into a render.RenderResponse
var _ PreImport = (*SayTemplate)(nil)

func (*SayTemplate) GetText(rt.Runtime) (g.Value, error) {
	panic("unexpected use of story method")
}

func (*SayTemplate) Execute(rt.Runtime) error {
	panic("unexpected use of story method")
}

// the import step transforms say response into a render.RenderResponse
var _ PreImport = (*SayResponse)(nil)

func (*SayResponse) GetText(rt.Runtime) (g.Value, error) {
	panic("unexpected use of story method")
}

func (*SayResponse) Execute(rt.Runtime) error {
	panic("unexpected use of story method")
}
