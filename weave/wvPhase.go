package weave

// Partitions the weave into stages;
// not entirely clear how much this is needed since missing resources retry
// helps with domain dependencies, plurals, and some amount of resource sorting
type Phase int

//go:generate stringer -type=Phase
const (
	RequireDependencies Phase = iota + 1 // domain
	RequirePlurals                       // kinds and kinds of kinds
	RequireAncestry                      // the actual kinds
	RequirePatterns                      // pattern definitions
	RequireDefaults
	RequireNouns
	RequireNames
	RequireRules
	RequireAll
)
