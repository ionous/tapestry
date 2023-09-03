package mdl

import (
	"git.sr.ht/~ionous/tapestry/lang"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
)

type PatternBuilder struct {
	Pattern
}

type Pattern struct {
	name, parent string
	fields       fieldSet
	rules        []Rule
	ruleOfs      int // tracks how many rules were successfully written ( for retru )
}

func (p Pattern) Copy(name string) Pattern {
	return Pattern{
		name:   name,
		parent: p.parent,
		fields: p.fields,
		rules:  p.rules,
	}
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
			name:   lang.Normalize(name),
			parent: parent,
		}}
}

//go:generate stringer -type=FieldType -linecomment
const (
	PatternParameters FieldType = iota // pattern parameters
	PatternResults
	PatternLocals
	NumFieldTypes
)

type FieldType int

// defers execution; so no return value.
func (b *PatternBuilder) AddField(ft FieldType, fn FieldInfo) {
	b.fields.fields[ft] = append(b.fields.fields[ft], fn)
}

func (b *PatternBuilder) AddLocal(fn FieldInfo) {
	b.AddField(PatternLocals, fn)
}
func (b *PatternBuilder) AddResult(fn FieldInfo) {
	b.AddField(PatternResults, fn)
}
func (b *PatternBuilder) AddParam(fn FieldInfo) {
	b.AddField(PatternParameters, fn)
}

// defers execution; so no return value.
// expects target class name to be normalized.
func (b *PatternBuilder) AppendRule(rank int, rule rt.Rule) {
	b.rules = append(b.rules, Rule{Rank: rank, Rule: rule})
}

func (pat *Pattern) writePattern(pen *Pen, create bool) (err error) {
	if cache, e := pat.fields.cache(pen); e != nil {
		err = e
	} else {
		var kid kindInfo
		if create {
			kid, err = pen.addKind(pat.name, pat.parent)
		} else {
			kid, err = pen.findRequiredKind(pat.name)
		}
		if err == nil {
			if e := pat.fields.writeFieldSet(pen, kid, cache); e != nil {
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
