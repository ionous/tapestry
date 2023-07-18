package mdl

import (
	"git.sr.ht/~ionous/tapestry/lang"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
	"git.sr.ht/~ionous/tapestry/weave/assert"
)

type PatternBuilder struct {
	Pattern
}

type Pattern struct {
	name, parent string
	fields       fieldSet
	rules        []rule
	ruleOfs      int
}

func (p *Pattern) Name() string {
	return p.name
}

func (p *Pattern) Parent() string {
	return p.parent
}

func (p *Pattern) NumFields(ft FieldType) int {
	return len(p.fields.fields[ft])
}

type rule struct {
	target string
	filter rt.BoolEval
	flags  assert.EventTiming
	prog   []rt.Execute
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
func (b *PatternBuilder) AddRule(target string, filter rt.BoolEval, flags assert.EventTiming, prog []rt.Execute) {
	b.rules = append(b.rules, rule{
		target: target,
		filter: filter,
		flags:  flags,
		prog:   prog,
	})
}

func (p *Pattern) writePattern(pen *Pen) (err error) {
	if kid, e := pen.addKind(p.name, p.parent); e != nil {
		err = e
	} else if e := p.fields.writeFieldSet(pen, kid); e != nil {
		err = e
	} else {
		for cnt := len(p.rules); p.ruleOfs < cnt; p.ruleOfs++ {
			rule := p.rules[p.ruleOfs]
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
				flags := fromTiming(rule.flags)
				if e := pen.addRule(kid, tgt, flags, filter, prog); e != nil {
					err = e
					break
				}
			}
		}
	}
	return
}

func fromTiming(timing assert.EventTiming) int {
	var part int
	always := timing&assert.RunAlways != 0
	if always {
		timing ^= assert.RunAlways
	}
	switch timing {
	case assert.Before:
		part = 0
	case assert.During:
		part = 1
	case assert.After:
		part = 2
	case assert.Later:
		part = 3
	}
	flags := part + int(rt.FirstPhase)
	if always {
		flags = -flags // marker for rules that need to always run (ex. counters "every third try" )
	}
	return flags
}
