package rt

type Jump int

//go:generate stringer -type=Jump
const (
	// transitions immediately.
	JumpNow Jump = iota
	// transitions after processing the current rule set.
	JumpNext
	// transition only after processing the entire phase.
	JumpLater
)

// Rule triggers a named series of statements when its filters and phase are satisfied.
// stopping before the action happens is considered a cancel.
type Rule struct {
	Name string
	// controls the style of transition:
	// when true, stops processing of any other rules; otherwise moves to the next phase.
	// a stop before the main rule is considered a canceled action;
	// and the result of the action in that case is false.
	Stop bool
	// controls when a transitions happens
	Jump Jump
	// whether the rule needs to be updated
	// before jumping to the next phase ( ex. for counters )
	Updates bool
	// a chain of if-else statements is considered an "optional" rule;
	// otherwise its considered a mandatory rule.
	// mandatory rules block other rules from running ( based on stop/stop )
	// optional rules only block other rules if some part of the if chain tested true.
	Exe []Execute
}
