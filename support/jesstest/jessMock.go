package jesstest

import (
	"fmt"
	"strings"

	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/dl/grammar"
	"git.sr.ht/~ionous/tapestry/dl/jess"
	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/support/inflect"
	"git.sr.ht/~ionous/tapestry/weave/mdl"
)

// implements Registrar to watch incoming calls.
// posted makes it more like a stub than a mock maybe? oh well.
type Mock struct {
	q                  jess.Query
	out                []string
	unique             map[string]int
	posted             [jess.PriorityCount][]jess.Process
	nounPool, nounPair map[string]string
}

func MakeMock(q jess.Query, nounPool, nounPair map[string]string) Mock {
	return Mock{q: q, nounPool: nounPool, nounPair: nounPair}
}

func (m *Mock) Generate(paragraph string) (ret []string, err error) {
	if e := jess.Generate(m.q, m, paragraph); e != nil {
		err = e
	} else if e := m.runPost(m.q); e != nil {
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

func (m *Mock) runPost(q jess.Query) (err error) {
Loop:
	for i := jess.Priority(0); i < jess.PriorityCount; i++ {
		// fix: can we really add new processes during post?
		// and if so, shouldnt mock panic on misorders?
		for len(m.posted[i]) > 0 {
			posted := m.posted[i]
			next, rest := posted[0], posted[1:]
			if e := next(q); e != nil {
				err = e
				break Loop
			} else {
				m.posted[i] = rest
			}
		}
	}
	return
}
func (m *Mock) AddNounKind(noun, kind string) (err error) {
	m.nounPool[noun] = noun
	// for weave, we'd add these blank kinds "objects"
	// we absorb it for cleaner tests;
	if len(kind) > 0 {
		if prev, ok := m.nounPool["$"+noun]; ok {
			// these hacks for testing sure are getting painful
			if prev == kind {
				err = fmt.Errorf("%w %s already declared as %s", mdl.Duplicate, noun, prev)
			} else {
				err = fmt.Errorf("%w %s already declared as %s", mdl.Conflict, noun, prev)
			}
		} else {
			m.out = append(m.out, "AddNounKind", noun, kind)
			m.nounPool["$"+noun] = kind
		}
	}
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
	m.nounPool[name] = noun
	lastNamedNoun = noun
	lastNamedSize = len(m.out)
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
func (m *Mock) AddNounTrait(name, trait string) (_ error) {
	m.out = append(m.out, "AddNounTrait", name, trait)
	return
}
func (m *Mock) AddNounValue(name, prop string, v rt.Assignment) (err error) {
	// prettify the output slightly
	var el any = v
	if t, ok := v.(*assign.FromText); ok {
		if _, ok := t.Value.(*literal.TextValue); ok {
			el = t.Value
		}
	} else if n, ok := v.(*assign.FromNumber); ok {
		if _, ok := n.Value.(literal.LiteralValue); ok {
			el = n.Value
		}
	}
	if str, e := Marshal(el); e != nil {
		err = e
	} else {
		m.out = append(m.out, "AddNounValue", name, prop, str)
	}
	return
}

func (m *Mock) AddNounPath(name string, parts []string, v literal.LiteralValue) (err error) {
	path := strings.Join(parts, ".")
	if str, e := Marshal(v); e != nil {
		err = e
	} else {
		m.out = append(m.out, "AddNounValue", name, path, str)
	}
	return
}
func (m *Mock) AddNounPair(rel, lhs, rhs string) (_ error) {
	if rel == "whereabouts" {
		m.nounPair[rhs] = lhs
	}
	m.out = append(m.out, "AddNounPair", lhs, rel, rhs)
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
func (m *Mock) AddFact(key string, partsAndValue ...string) (_ error) {
	m.out = append(m.out, "AddFact", key)
	m.out = append(m.out, partsAndValue...)
	return
}
func (m *Mock) Apply(verb jess.Macro, lhs, rhs []string) (_ error) {
	m.out = append(m.out, "ApplyMacro", verb.Name)
	m.out = append(m.out, lhs...)
	m.out = append(m.out, rhs...)
	if verb.Name == "contain" {
		for _, left := range lhs {
			for _, right := range rhs {
				m.nounPair[right] = left
			}
		}
	}
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

func (m *Mock) GetRelativeNouns(noun, relation string, primary bool) (ret []string, err error) {
	if relation != "whereabouts" || primary {
		err = fmt.Errorf("unexpected relation %v(primary: %v)", relation, primary)
	} else {
		if a, ok := m.nounPair[noun]; ok {
			ret = []string{a}
		}
	}
	return
}

func (m *Mock) GetOpposite(word string) (ret string, err error) {
	switch word {
	case "north":
		ret = "south"
	case "south":
		ret = "north"
	case "east":
		ret = "west"
	case "west":
		ret = "east"
	default:
		err = fmt.Errorf("unexpected opposition %q", word)
	}
	return
}
