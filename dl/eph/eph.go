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
	AspectPhase
	FieldPhase
	RelationPhase
	DefaultPhase
	NounPhase
	RelativePhase
	PatternPhase
	GrammarPhase
	TestPhase
	ReferencePhase
	//
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
