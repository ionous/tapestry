package rt

import "github.com/ionous/errutil"

const NoResult errutil.NoPanicError = "no result"

// Rule triggers a named series of statements when its filters and phase are satisfied.
type Rule struct {
	Name    string
	Filter  BoolEval
	Execute []Execute
	Updates bool
}
