package eph

import (
	"git.sr.ht/~ionous/iffy/dl/literal"
	"git.sr.ht/~ionous/iffy/jsn"
	"git.sr.ht/~ionous/iffy/jsn/cout"
	"github.com/ionous/errutil"
)

// implemented by individual commands
type Ephemera interface {
	// fix? remove catalog from the signature?
	Assemble(c *Catalog, d *Domain, at string) error
	Phase() Phase
}

type Phase int

//go:generate stringer -type=Phase
const (
	DomainPhase Phase = iota
	PluralPhase
	AncestryPhase
	AspectPhase // traits of kinds
	FieldPhase  // other properties of kinds
	NounPhase   // instances ( of kinds )
	ValuePhase
	RelativePhase // initial relations between nouns
	PatternPhase
	AliasPhase
	DirectivePhase // more grammar
	NumPhases
)

type PhaseActions map[Phase]PhaseAction

type PhaseAction struct {
	Flags PhaseFlags
	Do    func(d *Domain) (err error)
}

// wrapper for implementing Ephemera with free functions
type PhaseFunction struct {
	OnPhase Phase
	Do      func(*Catalog, *Domain, string) error
}

func (fn PhaseFunction) Phase() Phase { return fn.OnPhase }
func (fn PhaseFunction) Assemble(c *Catalog, d *Domain, at string) (err error) {
	return fn.Do(c, d, at)
}

type PhaseFlags struct {
	NoDuplicates bool
}

// shared generic marshal prog to text
func marshalout(cmd interface{}) (ret string, err error) {
	if cmd != nil {
		if m, ok := cmd.(jsn.Marshalee); !ok {
			err = errutil.New("can only marshal autogenerated types")
		} else {
			ret, err = cout.Marshal(m, literal.CompactEncoder)
		}
	}
	return
}

// we can use encode instead of marshal to get the raw unquoted values
// it works because everything here is a literal value.
// alt: give the literal interface a "get literal value" function.
func encodeLiteral(v literal.LiteralValue) (ret interface{}, err error) {
	if v != nil {
		if m, ok := v.(jsn.Marshalee); !ok {
			err = errutil.New("can only encode autogenerated types")
		} else if value, e := cout.Encode(m, literal.CompactEncoder); e != nil {
			err = e
		} else {
			ret = value
		}
	}
	return
}
