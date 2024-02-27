package jesstest

import (
	"fmt"
	"strings"

	"git.sr.ht/~ionous/tapestry/dl/jess"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/support/inflect"
	"git.sr.ht/~ionous/tapestry/weave/mdl"
)

// implements Registrar to watch incoming calls.
type Mock struct {
	out    []string
	unique map[string]int
}

func (m *Mock) AddFields(kind string, fields []mdl.FieldInfo) (_ error) {
	m.out = append(m.out, "AddFields", kind)
	for _, f := range fields {
		m.out = append(m.out, f.Name, f.Affinity.String(), f.Class)
	}
	return
}

func (m *Mock) AddKind(kind, ancestor string) (_ error) {
	m.out = append(m.out, "AddKind", kind, ancestor)
	return
}
func (m *Mock) AddKindTrait(kind, trait string) (_ error) {
	m.out = append(m.out, "AddKindTrait", kind, trait)
	return
}
func (m *Mock) AddPlural(many, one string) (_ error) {
	m.out = append(m.out, "AddPlural", many, one)
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
func (m *Mock) AddNounValue(name, prop string, v rt.Assignment) (err error) {
	if str, e := Marshal(v); e != nil {
		err = e
	} else {
		m.out = append(m.out, "AddNounValue", name, prop, str)
	}
	return
}
func (m *Mock) AddTraits(aspect string, traits []string) (err error) {
	if aspect != "color" && !strings.HasSuffix(aspect, " status") { // aspects are singular :/
		err = fmt.Errorf("unknown aspect %q", aspect)
	} else {
		m.out = append(m.out, "AddTraits", aspect)
		m.out = append(m.out, traits...)
	}
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
