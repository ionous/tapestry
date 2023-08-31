package grammar

import (
	"git.sr.ht/~ionous/tapestry/parser"
)

// GrammarMaker -
type GrammarMaker interface{ MakeGrammar() GrammarMaker }

func (op *Directive) MakeGrammar() GrammarMaker { return op }

// acts as AllOf{Words, ... }
func (op *Directive) MakeScanners() (ret parser.Scanner) {
	return makeSequence(op.Series)
}
