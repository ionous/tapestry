package mdl

// Partitions the weave into stages;
// not entirely clear how much this is needed since missing resources retry
// helps with domain dependencies, plurals, and some amount of resource sorting
type Phase int

//go:generate stringer -type=Phase
const (
	DependencyPhase Phase = iota + 1
	LanguagePhase         // definitions of words
	AncestryPhase         // kinds and their derivation
	PropertyPhase         // the members of kinds
	MappingPhase          // match phrases that might otherwise treat directions as names.
	NounPhase             // generate explicit nouns
	MacroPhase            // tbd: merge with nouns?
	ConnectionPhase       // pair up nouns, sometimes implying new nouns or specific kinds.
	FallbackPhase         // generate kinds for nouns that didn't derive during connections
	ValuePhase            // apply any collected values ( tbd: merge with rules? )
	RulePhase
	FinalPhase
	NumPhases
)
