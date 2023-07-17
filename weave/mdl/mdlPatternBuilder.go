package mdl

import (
	"git.sr.ht/~ionous/tapestry/lang"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
	"git.sr.ht/~ionous/tapestry/weave/assert"
	"git.sr.ht/~ionous/tapestry/weave/res"
)

type PatternBuilder struct {
	search GetRequiredKind
	Pattern
}

type Pattern struct {
	Kind
	rules   []rule
	ruleOfs int
}

func (p *Pattern) Name() string {
	return p.name
}

func (p *Pattern) Parent() string {
	return p.parent
}

func (p *Pattern) NumFields(ft FieldType) int {
	return len(p.fs.fields[ft])
}

type rule struct {
	target string
	filter rt.BoolEval
	flags  assert.EventTiming
	prog   []rt.Execute
}

type classCache map[string]res.Result

func (b *classCache) addClass(kinds GetRequiredKind, cls string) {
	if len(cls) > 0 && (*b)[cls] == nil {
		if (*b) == nil {
			(*b) = make(map[string]res.Result)
		}
		(*b)[cls] = kinds.GetRequiredKind(cls)
	}
}

func (b *classCache) getClass(cls string) (ret KindInfo, err error) {
	if r, ok := (*b)[cls]; ok {
		if v, e := r.Resolve(); e != nil {
			err = e
		} else {
			ret = v.(KindInfo)
		}
	}
	return
}

type GetRequiredKind interface {
	GetRequiredKind(kind string) res.Result
}

func NewPatternBuilder(w GetRequiredKind, name string) *PatternBuilder {
	return NewPatternSubtype(w, name, kindsOf.Pattern)
}

func NewPatternSubtype(w GetRequiredKind, name string, parent kindsOf.Kinds) *PatternBuilder {
	if (parent & kindsOf.Pattern) == 0 {
		panic("subtype not a pattern")
	}
	return &PatternBuilder{
		search: w,
		Pattern: Pattern{
			Kind: Kind{
				// tbd: feels like it'd be best to have spec flag names that need normalization,
				// and convert all the names at load time ( probably storing the original somewhere )
				// ( ex. store the normalized names in the meta data )
				name:   lang.Normalize(name),
				parent: parent.String(),
			},
		}}
}

// defers execution; so no return value.
func (b *PatternBuilder) AddField(ft FieldType, fn FieldInfo) {
	b.classes.addClass(b.search, fn.Class)
	b.fs.fields[ft] = append(b.fs.fields[ft], fn)
}

// defers execution; so no return value.
// expects target class name to be normalized.
func (b *PatternBuilder) AddRule(target string, filter rt.BoolEval, flags assert.EventTiming, prog []rt.Execute) {
	b.classes.addClass(b.search, target)
	b.rules = append(b.rules, rule{
		target: target,
		filter: filter,
		flags:  flags,
		prog:   prog,
	})
}

func (p *Pattern) write(m *Pen) (ret KindInfo, err error) {
	// if kid, e := p.Kind.write(m); e != nil {
	// 	err = e
	// } else {
	// 	for cnt := len(p.rules); p.ruleOfs < cnt; p.ruleOfs++ {
	// 		rule := p.rules[p.ruleOfs]
	// 		if tgt, e := p.classes.getClass(rule.target); e != nil {
	// 			err = e
	// 			break
	// 		} else {
	// 			flags := fromTiming(rule.flags)
	// 			if e := m.AddRuleById(kid, tgt, flags, rule.filter, rule.prog); e != nil {
	// 				err = e
	// 				break
	// 			}
	// 		}
	// 	}
	// 	if err == nil {
	// 		ret = kid
	// 	}
	// }
	// return
	panic("not implemented")
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
