package eph

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
	AspectPhase   // traits of kinds
	FieldPhase    // other properties of kinds
	DefaultPhase  // default values of fields
	NounPhase     // instances ( of kinds )
	RelativePhase // initial relations between nouns
	PatternPhase
	GrammarPhase
	DirectivePhase // more grammar
	TestPhase
	ReferencePhase // tdb: names used by things already built;
	/*          */ // might be redundant if patterns, etc. are looking up the things they need
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
