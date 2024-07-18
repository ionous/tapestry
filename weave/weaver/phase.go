package weaver

// Partitions the weave to order dependencies between types.
// ( ex. Nouns depend on kinds, and never the other way around. )
type Phase int

//go:generate stringer -type=Phase
const (
	// the zeroth phase is never explicitly scheduled.
	LanguagePhase   Phase = iota + 1 // definitions of words
	AncestryPhase                    // kinds and their derivation
	PropertyPhase                    // the members of kinds ( after ancestry because fields depend on kind )
	NounPhase                        // generate explicit nouns
	VerbPhase                        // apply existing verbs
	ConnectionPhase                  // pairings and map connections
	FallbackPhase                    // generate kinds for nouns that didn't derive during connections
	ValuePhase                       // apply any collected values
	NextPhase                        // any sub domains
	NumPhases
)
