package grammar

import (
	"git.sr.ht/~ionous/tapestry/parser"
)

// acts as AllOf{Words, ... }
func (op *Directive) MakeScanners() (ret parser.Scanner) {
	return makeSequence(op.Series)
}
