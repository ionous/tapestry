package testutil

import (
	"git.sr.ht/~ionous/iffy/affine"
	"git.sr.ht/~ionous/iffy/rt"
	g "git.sr.ht/~ionous/iffy/rt/generic"
	"git.sr.ht/~ionous/iffy/rt/writer"
)

// PanicRuntime implements Runtime throwing a panic for every method
type PanicRuntime struct{}

var _ rt.Runtime = (*PanicRuntime)(nil)

func (PanicRuntime) ActivateDomain(name string, enable bool) {
	panic("Runtime panic")
}
func (PanicRuntime) GetKindByName(string) (*g.Kind, error) {
	panic("Runtime panic")
}
func (PanicRuntime) Call(name string, aff affine.Affinity, args []rt.Arg) (ret g.Value, err error) {
	panic("Runtime panic")
}
func (PanicRuntime) RelateTo(a, b, relation string) error {
	panic("Runtime panic")
}
func (PanicRuntime) RelativesOf(a, relation string) ([]string, error) {
	panic("Runtime panic")
}
func (PanicRuntime) ReciprocalsOf(a, relation string) ([]string, error) {
	panic("Runtime panic")
}
func (PanicRuntime) GetField(target, field string) (g.Value, error) {
	panic("Runtime panic")
}
func (PanicRuntime) SetField(target, field string, v g.Value) error {
	panic("Runtime panic")
}
func (PanicRuntime) Writer() writer.Output {
	panic("Runtime panic")
}
func (PanicRuntime) SetWriter(writer.Output) writer.Output {
	panic("Runtime panic")
}
func (PanicRuntime) PushScope(rt.Scope) {
	panic("Runtime panic")
}
func (PanicRuntime) PopScope() {
	panic("Runtime panic")
}
func (PanicRuntime) ReplaceScope(rt.Scope) rt.Scope {
	panic("Runtime panic")
}
func (PanicRuntime) Random(inclusiveMin, exclusiveMax int) int {
	panic("Runtime panic")
}
func (PanicRuntime) PluralOf(single string) string {
	panic("Runtime panic")
}
func (PanicRuntime) SingularOf(plural string) string {
	panic("Runtime panic")
}
