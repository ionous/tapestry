package weaver

import (
	"git.sr.ht/~ionous/tapestry/dl/grammar"
	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/weave/mdl"
)

type PanicWeaves struct{}

var _ Weaves = (*PanicWeaves)(nil)

// AddAspectTraits implements Weaves.
func (p *PanicWeaves) AddAspectTraits(aspect string, traits []string) error {
	panic("unimplemented")
}

// AddCheck implements Weaves.
func (p *PanicWeaves) AddCheck(name string, value literal.LiteralValue, prog []rt.Execute) error {
	panic("unimplemented")
}

// AddFact implements Weaves.
func (p *PanicWeaves) AddFact(key string, partsAndValue ...string) error {
	panic("unimplemented")
}

// AddGrammar implements Weaves.
func (p *PanicWeaves) AddGrammar(name string, prog *grammar.Directive) error {
	panic("unimplemented")
}

// AddKind implements Weaves.
func (p *PanicWeaves) AddKind(kind string, ancestor string) error {
	panic("unimplemented")
}

// AddKindFields implements Weaves.
func (p *PanicWeaves) AddKindFields(kind string, fields []mdl.FieldInfo) error {
	panic("unimplemented")
}

// AddKindTrait implements Weaves.
func (p *PanicWeaves) AddKindTrait(kind string, trait string) error {
	panic("unimplemented")
}

// AddNounKind implements Weaves.
func (p *PanicWeaves) AddNounKind(noun string, kind string) error {
	panic("unimplemented")
}

// AddNounName implements Weaves.
func (p *PanicWeaves) AddNounName(noun string, name string, rank int) error {
	panic("unimplemented")
}

// AddNounPair implements Weaves.
func (p *PanicWeaves) AddNounPair(rel string, noun string, otherNoun string) error {
	panic("unimplemented")
}

// AddNounPath implements Weaves.
func (p *PanicWeaves) AddNounPath(noun string, path []string, val literal.LiteralValue) error {
	panic("unimplemented")
}

// AddNounTrait implements Weaves.
func (p *PanicWeaves) AddNounTrait(noun string, trait string) error {
	panic("unimplemented")
}

// AddNounValue implements Weaves.
func (p *PanicWeaves) AddNounValue(noun string, prop string, val rt.Assignment) error {
	panic("unimplemented")
}

// AddOpposite implements Weaves.
func (p *PanicWeaves) AddOpposite(a string, b string) error {
	panic("unimplemented")
}

// AddPattern implements Weaves.
func (p *PanicWeaves) AddPattern(mdl.Pattern) error {
	panic("unimplemented")
}

// AddPlural implements Weaves.
func (p *PanicWeaves) AddPlural(many string, one string) error {
	panic("unimplemented")
}

// AddRelation implements Weaves.
func (p *PanicWeaves) AddRelation(name string, oneKind string, otherKind string, amany bool, bmany bool) error {
	panic("unimplemented")
}

// ExtendPattern implements Weaves.
func (p *PanicWeaves) ExtendPattern(mdl.Pattern) error {
	panic("unimplemented")
}

// GenerateUniqueName implements Weaves.
func (p *PanicWeaves) GenerateUniqueName(category string) string {
	panic("unimplemented")
}
