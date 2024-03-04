package jesstest

import (
	"fmt"
	"strings"

	"git.sr.ht/~ionous/tapestry/dl/grammar"
	"git.sr.ht/~ionous/tapestry/dl/jess"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/support/inflect"
	"git.sr.ht/~ionous/tapestry/weave/mdl"
)

// implements Registrar to watch incoming calls.
// posted makes it more like a stub than a mock maybe? oh well.
type Mock struct {
	q        jess.Query
	out      []string
	unique   map[string]int
	posted   [jess.PriorityCount][]jess.Process
	nounPool map[string]bool
}

func MakeMock(q jess.Query, nounPool map[string]bool) Mock {
	return Mock{q: q, nounPool: nounPool}
}

func (m *Mock) Generate(str string) (ret []string, err error) {
	if e := jess.Generate(m.q, m, str); e != nil {
		err = e
	} else if e := m.RunPost(m.q); e != nil {
		err = e
	} else {
		ret = m.out
	}
	return
}

func (m *Mock) PostProcess(i jess.Priority, p jess.Process) (_ error) {
	m.posted[i] = append(m.posted[i], p)
	return
}

func (m *Mock) RunPost(q jess.Query) (err error) {
	for _, posted := range m.posted {
		// fix: can we really add new processes during post?
		// and if so, shouldnt mock panic on misorders?
		for len(posted) > 0 {
			next, rest := posted[0], posted[1:]
			if e := next(q, m); e != nil {
				err = e
				break
			} else {
				posted = rest
			}
		}
	}
	return
}

func (m *Mock) AddFields(kind string, fields []mdl.FieldInfo) (_ error) {
	m.out = append(m.out, "AddFields", kind)
	for _, f := range fields {
		m.out = append(m.out, f.Name, f.Affinity.String(), f.Class)
	}
	return
}

func (m *Mock) AddGrammar(name string, prog *grammar.Directive) (err error) {
	m.out = append(m.out, "AddGrammar", name)
	// we know the top level format of the grammar;
	// so break it open to make the test output easier to read
	action := prog.Series[1].(*grammar.Action)
	series := prog.Series[0].(*grammar.Sequence)
	for _, n := range append(series.Series, action) {
		if str, e := Marshal(n); e != nil {
			err = e
			break
		} else {
			m.out = append(m.out, str)
		}
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
func (m *Mock) AddNounKind(noun, kind string) (_ error) {
	m.nounPool[noun] = true
	m.out = append(m.out, "AddNounKind", noun, kind)
	return
}

var lastNamedNoun string
var lastNamedSize int

// slightly limit the name spew; name generation gets tested elsewhere
func (m *Mock) AddNounName(noun, name string, r int) (_ error) {
	if lastNamedSize != len(m.out) {
		lastNamedNoun = ""
	}
	if r < 0 {
		m.out = append(m.out, "AddNounAlias", noun, name)
	} else if lastNamedNoun != noun || r < 0 {
		m.out = append(m.out, "AddNounName", noun, name)
	}

	lastNamedNoun = noun
	lastNamedSize = len(m.out)
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
