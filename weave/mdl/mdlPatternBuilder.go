package mdl

import (
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
	"git.sr.ht/~ionous/tapestry/support/inflect"
)

type PatternBuilder struct {
	Pattern
}

type Rule struct {
	rt.Rule
	Rank int
}

// public snapshot of the desired pattern
type Pattern struct {
	name, parent string
	fields       patternFields
	rules        []Rule
	ruleOfs      int // number of successfully written rules
}

func (p *Pattern) Name() string {
	return p.name
}

func (p *Pattern) Parent() string {
	return p.parent
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
	const i = patternParameters
	b.fields[i] = append(b.fields[i], fs...)
}

// defers execution; so no return value.
func (b *PatternBuilder) AddResult(f FieldInfo) {
	const i = patternResults
	b.fields[i] = append(b.fields[i], f)
}

// defers execution; so no return value.
func (b *PatternBuilder) AddLocals(fs []FieldInfo) {
	const i = patternLocals
	b.fields[i] = append(b.fields[i], fs...)
}

// defers execution; so no return value.
// expects target class name to be normalized.
func (b *PatternBuilder) AppendRule(rank int, rule rt.Rule) {
	b.rules = append(b.rules, Rule{Rank: rank, Rule: rule})
}

func (pat *Pattern) writePattern(pen *Pen, create bool) (err error) {
	// make sure all of the pattern dependencies are known before trying to generate anything.
	// note: patterns marked "create" will block out the other kinds
	var kind kindInfo
	if !create {
		kind, err = pen.findRequiredKind(pat.name)
	}
	if err == nil {
		if cache, e := pat.fields.precache(pen); e != nil {
			err = e
		} else {
			if create {
				kind, err = pen.addKind(pat.name, pat.parent)
			}
			if err == nil {
				err = pat.writeFields(pen, kind, cache)
			}
		}
	}
	return
}

func (pat *Pattern) writeFields(pen *Pen, kind kindInfo, cache fieldCache) (err error) {
	var blank FieldInfo
	var handlers = [3]fieldHandler{pen.addParameter, pen.addResult, pen.addField}
	for i, fields := range pat.fields {
		fieldHandler := handlers[i]
		if len(fields) == 1 && fields[0] == blank {
			// the Nothing type generates a blank field info
			// fix; should probably be an error if nothing is used for locals
		} else {
			fs := fieldSet{kind, fields, cache}
			if e := fs.addFields(pen, fieldHandler); e != nil {
				err = e
			} else {
				for cnt := len(pat.rules); pat.ruleOfs < cnt; pat.ruleOfs++ {
					rule := pat.rules[pat.ruleOfs]
					if prog, e := marshalprog(rule.Exe); e != nil {
						err = e
					} else if e := pen.addRule(kind,
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

type patternFields [numPatternParts][]FieldInfo

// the order here reflects the order in the db
const (
	patternParameters int = iota // pattern parameters
	patternResults
	patternLocals
	numPatternParts
)

func (parts *patternFields) precache(pen *Pen) (cache fieldCache, err error) {
	for _, fields := range *parts {
		if e := cache.precache(pen, fields); e != nil {
			err = e
			break
		}
	}
	return
}
