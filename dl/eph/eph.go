package eph

import (
	"git.sr.ht/~ionous/iffy/dl/literal"
	"git.sr.ht/~ionous/iffy/jsn"
	"git.sr.ht/~ionous/iffy/jsn/cout"
	"github.com/ionous/errutil"
)

// database/sql like interface
type Writer interface {
	Write(q string, args ...interface{}) error
}

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
	AspectPhase    // traits of kinds
	FieldPhase     // collect the properties of kinds
	PostFieldPhase // actually assemble those fields
	NounPhase      // instances ( of kinds )
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
