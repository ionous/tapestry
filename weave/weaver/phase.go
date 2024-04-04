package weaver

// Partitions the weave into stages;
// not entirely clear how much this is needed since missing resources retry
// helps with domain dependencies, plurals, and some amount of resource sorting
type Phase int

//go:generate stringer -type=Phase
const (
	DependencyPhase Phase = iota + 1
	LanguagePhase         // definitions of words
	AncestryPhase         // kinds and their derivation
	PropertyPhase         // the members of kinds ( after ancestry because fields depend on kind )
	NounPhase             // generate explicit nouns
	VerbPhase             // apply existing verbs
	ConnectionPhase       // pairings and map connections
	FallbackPhase         // generate kinds for nouns that didn't derive during connections
	ValuePhase            // apply any collected values
	NextPhase             // any sub domains
	NumPhases
)
