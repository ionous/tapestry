package testutil

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/rt"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"git.sr.ht/~ionous/tapestry/rt/writer"
)

// PanicRuntime implements Runtime throwing a panic for every method
type PanicRuntime struct{}

var _ rt.Runtime = (*PanicRuntime)(nil)

func (PanicRuntime) ActivateDomain(name string) (string, error) {
	panic("Runtime panic")
}
func (PanicRuntime) GetKindByName(string) (*g.Kind, error) {
	panic("Runtime panic")
}
func (PanicRuntime) GetRules(pattern, target string, pflags *rt.Flags) (ret []rt.Rule, err error) {
	panic("Runtime panic")
}
func (PanicRuntime) Call(name string, aff affine.Affinity, args []rt.Arg) (ret g.Value, err error) {
	panic("Runtime panic")
}
func (PanicRuntime) Send(pat string, up []string, args []rt.Arg) (ret g.Value, err error) {
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
func (PanicRuntime) ReplaceScope(scope rt.Scope, init bool) (ret rt.Scope, err error) {
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
