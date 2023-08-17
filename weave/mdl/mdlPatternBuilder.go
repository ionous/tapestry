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
	rules        []rule
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

type rule struct {
	target string
	filter rt.BoolEval
	prog   []rt.Execute
	appends,
	updates,
	terminates bool
}

func NewPatternBuilder(name string) *PatternBuilder {
	return NewPatternSubtype(name, kindsOf.Pattern)
}

func NewPatternSubtype(name string, parent kindsOf.Kinds) *PatternBuilder {
	if (parent & kindsOf.Pattern) == 0 {
		panic("subtype not a pattern")
	}
	return &PatternBuilder{
		Pattern: Pattern{
			// tbd: feels like it'd be best to have spec flag names that need normalization,
			// and convert all the names at load time ( probably storing the original somewhere )
			// ( ex. store the normalized names in the meta data )
			name:   lang.Normalize(name),
			parent: parent.String(),
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
func (b *PatternBuilder) AddRule(target string, filter rt.BoolEval, updates bool, prog []rt.Execute) {
	b.rules = append(b.rules, rule{
		target:  target,
		filter:  filter,
		updates: updates,
		prog:    prog,
	})
}

func (b *PatternBuilder) AddNewRule(name string, appends, updates, terminates bool, prog []rt.Execute) {
	b.rules = append(b.rules, rule{
		// fix: name,
		appends:    appends,
		updates:    updates,
		terminates: terminates,
		prog:       prog,
	})
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
					if tgt, e := pen.findOptionalKind(rule.target); e != nil {
						err = e
						break
					} else if filter, e := marshalout(rule.filter); e != nil {
						err = e
						break
					} else if prog, e := marshalprog(rule.prog); e != nil {
						err = e
						break
					} else {
						if e := pen.addRule(kid, tgt, 0, rule.appends, rule.updates, rule.terminates, filter, prog); e != nil {
							err = e
							break
						}
					}
				}
			}
		}
	}
	return
}
