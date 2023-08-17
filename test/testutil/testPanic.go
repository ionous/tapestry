package testutil

import (
	"io"

	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/rt"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
)

// PanicRuntime implements Runtime throwing a panic for every method
type PanicRuntime struct{}

var _ rt.Runtime = (*PanicRuntime)(nil)

func (PanicRuntime) ActivateDomain(name string) (string, error) {
	panic("Runtime panic")
}
func (PanicRuntime) Call(string, affine.Affinity, []string, []g.Value) (g.Value, error) {
	panic("Runtime panic")
}
func (PanicRuntime) GetField(target, field string) (g.Value, error) {
	panic("Runtime panic")
}
func (PanicRuntime) GetKindByName(string) (*g.Kind, error) {
	panic("Runtime panic")
}
func (PanicRuntime) GetRules(pattern, target string) ([]rt.Rule, error) {
	panic("Runtime panic")
}
func (PanicRuntime) OppositeOf(string) string {
	panic("Runtime panic")
}
func (PanicRuntime) PluralOf(single string) string {
	panic("Runtime panic")
}
func (PanicRuntime) PopScope() {
	panic("Runtime panic")
}
func (PanicRuntime) PushScope(rt.Scope) {
	panic("Runtime panic")
}
func (PanicRuntime) RelateTo(a, b, relation string) error {
	panic("Runtime panic")
}
func (PanicRuntime) RelativesOf(a, relation string) (g.Value, error) {
	panic("Runtime panic")
}
func (PanicRuntime) ReciprocalsOf(a, relation string) (g.Value, error) {
	panic("Runtime panic")
}
func (PanicRuntime) Send(*g.Record, []string) (g.Value, error) {
	panic("Runtime panic")
}
func (PanicRuntime) SetField(target, field string, v g.Value) error {
	panic("Runtime panic")
}
func (PanicRuntime) SetWriter(io.Writer) io.Writer {
	panic("Runtime panic")
}
func (PanicRuntime) SingularOf(plural string) string {
	panic("Runtime panic")
}
func (PanicRuntime) Writer() io.Writer {
	panic("Runtime panic")
}
func (PanicRuntime) Random(inclusiveMin, exclusiveMax int) int {
	panic("Runtime panic")
}
