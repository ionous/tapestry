package jesstest

import (
	"errors"
	"fmt"
	"strings"

	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/dl/grammar"
	"git.sr.ht/~ionous/tapestry/dl/jess"
	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/weave/mdl"
	"git.sr.ht/~ionous/tapestry/weave/weaver"
)

// implements Registrar to watch incoming calls.
// posted makes it more like a stub than a mock maybe? oh well.
type Mock struct {
	weaver.PanicWeaves
	q                   jess.Query
	out                 []string
	unique              map[string]int
	nounPool, nounPairs map[string]string
	ProcessingList
	jessRt jessRt
}

func MakeMock(q jess.Query, nouns map[string]string) Mock {
	pairs := make(map[string]string)
	return Mock{
		q:         q,
		nounPool:  nouns,
		nounPairs: pairs,
		jessRt:    jessRt{nounPairs: pairs, verbs: KnownVerbs},
	}
}

func (m *Mock) Generate(str string, val rt.Assignment) (ret []string, err error) {
	if e := m.generate(str, val); e != nil {
		err = e
	} else {
		ret = m.out
	}
	return
}

func (m *Mock) generate(str string, val rt.Assignment) (err error) {
	if p, e := jess.NewParagraph(str, val); e != nil {
		err = e
	} else {
		for z := weaver.Phase(0); z < weaver.NumPhases; z++ {
			if _, e := p.Generate(z, m.q, m); e != nil {
				err = e // match, and schedule callbacks for (later) phases
				break
			} else {
				// update callbacks for the current phase
				// in story; these would be intermixed with other scheduled elements for the phase
				cnt, e := m.UpdatePhase(z, m, &m.jessRt)
				missing := errors.Is(e, weaver.Missing)
				if (e != nil && !missing) || (missing && cnt == 0) {
					err = e
					break
				}
			}
		}
	}
	return
}

func (m *Mock) AddNounKind(noun, kind string) (err error) {
	m.nounPool[noun] = noun
	if len(kind) == 0 {
		err = fmt.Errorf("expected a valid kind for %q", noun)
	} else if kind != jess.Objects {
		// ^ absorb this for cleaner test ouptut
		if prev, exists := m.nounPool["$"+noun]; !exists {
			m.addNounKind(noun, kind)
		} else if prev == kind {
			err = fmt.Errorf("%w %s already declared as %s", weaver.Duplicate, noun, prev)
		} else if prev == "things" && thingLike(kind) {
			m.addNounKind(noun, kind)
		} else if kind != "things" || !thingLike(prev) {
			err = fmt.Errorf("%w %s already declared as %s", weaver.Conflict, noun, prev)
		}
	}
	return
}

func thingLike(k string) (okay bool) {
	switch k {
	case "doors", "containers", "supporters":
		okay = true
	}
	return
}

func (m *Mock) addNounKind(noun, kind string) {
	m.out = append(m.out, "AddNounKind:", noun, kind)
	m.nounPool["$"+noun] = kind
}

var lastNamedNoun string
var lastNamedSize int

// slightly limit the name spew; name generation gets tested elsewhere
func (m *Mock) AddNounName(noun, name string, r int) (_ error) {
	if lastNamedSize != len(m.out) {
		lastNamedNoun = ""
	}
	if r < 0 {
		m.out = append(m.out, "AddNounAlias:", noun, name)
	} else if lastNamedNoun != noun || r < 0 {
		m.out = append(m.out, "AddNounName:", noun, name)
	}
	m.nounPool[name] = noun
	lastNamedNoun = noun
	lastNamedSize = len(m.out)
	return
}

func (m *Mock) AddKindFields(kind string, fields []mdl.FieldInfo) (_ error) {
	m.out = append(m.out, "AddKindFields:", kind)
	for _, f := range fields {
		m.out = append(m.out, f.Name, f.Affinity.String(), f.Class)
	}
	return
}

func (m *Mock) AddGrammar(name string, prog *grammar.Directive) (err error) {
	m.out = append(m.out, "AddGrammar:", name)
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
	m.out = append(m.out, "AddKind:", kind, ancestor)
	return
}
func (m *Mock) AddKindTrait(kind, trait string) (_ error) {
	m.out = append(m.out, "AddKindTrait:", kind, trait)
	return
}
func (m *Mock) AddPlural(many, one string) (_ error) {
	m.out = append(m.out, "AddPlural:", many, one)
	return
}
func (m *Mock) AddNounTrait(name, trait string) (_ error) {
	m.out = append(m.out, "AddNounTrait:", name, trait)
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
		m.out = append(m.out, "AddNounValue:", name, prop, str)
	}
	return
}

func (m *Mock) AddNounPath(name string, parts []string, v literal.LiteralValue) (err error) {
	path := strings.Join(parts, ".")
	if str, e := Marshal(v); e != nil {
		err = e
	} else {
		m.out = append(m.out, "AddNounValue:", name, path, str)
	}
	return
}
func (m *Mock) AddNounPair(rel, lhs, rhs string) (_ error) {
	if rel == "whereabouts" {
		m.nounPairs[rhs] = lhs
	}
	m.out = append(m.out, "AddNounPair:", rel, lhs, rhs)
	return
}
func (m *Mock) AddAspectTraits(aspect string, traits []string) (err error) {
	if aspect != "color" && !strings.HasSuffix(aspect, " status") { // aspects are singular :/
		err = fmt.Errorf("unknown aspect %q", aspect)
	} else {
		m.out = append(m.out, "AddAspectTraits:", aspect)
		m.out = append(m.out, traits...)
	}
	return
}

func (m *Mock) ExtendPattern(p mdl.Pattern) (_ error) {
	m.out = append(m.out, "ExtendPattern:", p.Name(),
		fmt.Sprintf("rule count: %d", len(p.Rules())))
	return
}

// mock assumes all facts valid and new
func (m *Mock) AddFact(key string, partsAndValue ...string) (_ error) {
	m.out = append(m.out, "AddFact:", key)
	m.out = append(m.out, partsAndValue...)
	return
}

func (m *Mock) GenerateUniqueName(category string) string {
	if m.unique == nil {
		m.unique = make(map[string]int)
	}
	next := m.unique[category] + 1
	m.unique[category] = next
	return fmt.Sprintf("%s-%d", category, next)
}
