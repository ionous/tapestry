package mdl

import (
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
	"git.sr.ht/~ionous/tapestry/support/inflect"
)

type PatternBuilder struct {
	Pattern
}

// public snapshot of the desired pattern
type Pattern struct {
	name, parent            string
	params, results, locals []FieldInfo // the order here reflects the order in the db
	rules                   []Rule
	ruleOfs                 int // number of successfully written rules
}

func (p *Pattern) Name() string {
	return p.name
}

func (p *Pattern) Parent() string {
	return p.parent
}

type Rule struct {
	rt.Rule
	Rank int
}

func NewPatternBuilder(name string) *PatternBuilder {
	return NewPatternSubtype(name, kindsOf.Pattern.String())
}

func NewPatternSubtype(name string, parent string) *PatternBuilder {
	return &PatternBuilder{
		Pattern: Pattern{
			// tbd: feels like it'd be best to have spec flag names that need normalization,
			// and convert all the names at load time ( probably storing the original somewhere )
			// ( ex. store the normalized names in the meta data )
			name:   inflect.Normalize(name),
			parent: parent,
		}}
}

// defers execution; so no return value.
func (b *PatternBuilder) AddParams(fs []FieldInfo) {
	b.params = append(b.params, fs...)
}

// defers execution; so no return value.
func (b *PatternBuilder) AddResult(f FieldInfo) {
	b.results = append(b.results, f)
}

// defers execution; so no return value.
func (b *PatternBuilder) AddLocals(fs []FieldInfo) {
	b.locals = append(b.locals, fs...)
}

// defers execution; so no return value.
// expects target class name to be normalized.
func (b *PatternBuilder) AppendRule(rank int, rule rt.Rule) {
	b.rules = append(b.rules, Rule{Rank: rank, Rule: rule})
}

func (pat *Pattern) writePattern(pen *Pen, create bool) (err error) {
	var cache fieldCache
	var parts = [3]struct {
		fields []FieldInfo
		fieldHandler
	}{{
		pat.params,
		pen.addParameter,
	}, {
		pat.results,
		pen.addResult,
	}, {
		pat.locals,
		pen.addField,
	},
	}
	var kid kindInfo
	if create {
		kid, err = pen.addKind(pat.name, pat.parent)
	} else {
		kid, err = pen.findRequiredKind(pat.name)
	}
	if err == nil {
		var blank FieldInfo
		for _, p := range parts {
			if len(p.fields) == 1 && p.fields[0] == blank {
				// the Nothing type generates a blank field info
				// fix; should probably be an error if nothing is used for locals
			} else if e := cache.writeFields(pen, kid, p.fields, p.fieldHandler); e != nil {
				err = e
			} else {
				for cnt := len(pat.rules); pat.ruleOfs < cnt; pat.ruleOfs++ {
					rule := pat.rules[pat.ruleOfs]
					if prog, e := marshalprog(rule.Exe); e != nil {
						err = e
					} else if e := pen.addRule(kid,
						rule.Name, rule.Rank,
						rule.Stop, int(rule.Jump),
						rule.Updates, prog); e != nil {
						err = e
						break
					}
				}
			}
		}
	}
	return
}
