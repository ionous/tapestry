package jesstest

import (
	"fmt"

	"git.sr.ht/~ionous/tapestry/dl/jess"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/support/inflect"
)

// implements Registrar to watch incoming calls.
type Mock struct {
	out    []string
	unique map[string]int
}

func (m *Mock) AddKind(kind, ancestor string) (_ error) {
	m.out = append(m.out, "AddKind", kind, ancestor)
	return
}
func (m *Mock) AddKindTrait(kind, trait string) (_ error) {
	m.out = append(m.out, "AddKindTrait", kind, trait)
	return
}
func (m *Mock) AddNoun(short, long, kind string) (_ error) {
	m.out = append(m.out, "AddNoun", short, long, kind)
	return
}
func (m *Mock) AddNounAlias(noun, name string, _ int) (_ error) {
	m.out = append(m.out, "AddNounAlias", noun, name)
	return
}
func (m *Mock) AddNounTrait(name, trait string) (_ error) {
	m.out = append(m.out, "AddNounTrait", name, trait)
	return
}
func (m *Mock) AddNounValue(name, prop string, _ rt.Assignment) (_ error) {
	m.out = append(m.out, "AddNounValue", name, prop)
	return
}
func (m *Mock) Apply(verb jess.Macro, lhs, rhs []string) (_ error) {
	m.out = append(m.out, "ApplyMacro", verb.Name)
	m.out = append(m.out, lhs...)
	m.out = append(m.out, rhs...)
	return
}
func (m *Mock) GetClosestNoun(name string) (string, error) {
	return name, nil
}
func (m *Mock) GetExactNoun(name string) (string, error) {
	return name, nil
}
func (m *Mock) GetPlural(word string) string {
	return inflect.Pluralize(word)
}
func (m *Mock) GetSingular(word string) string {
	return inflect.Singularize(word)
}
func (m *Mock) GetUniqueName(category string) string {
	if m.unique == nil {
		m.unique = make(map[string]int)
	}
	next := m.unique[category] + 1
	m.unique[category] = next
	return fmt.Sprintf("%s-%d", category, next)
}
