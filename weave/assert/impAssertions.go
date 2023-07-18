package assert

type EventTiming int

const (
	//
	DefaultTiming EventTiming = iota
	Before
	During
	After
	Later     // 100
	RunAlways = 1 << 3
)

type Phase int

//go:generate stringer -type=Phase
const (
	RequireDependencies Phase = iota + 1
	RequirePlurals
	RequireDeterminers // kinds and kinds of kinds
	RequireAncestry    // parameters precede results
	RequireParameters  // results precede normal fields
	RequireResults     // properties of kinds and records
	RequireFields      // initial values of fields
	RequireDefaults
	RequireNouns
	RequireNames
	RequireRelatives
	RequireRules
	RequireAll
)
