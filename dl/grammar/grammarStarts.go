package grammar

import (
	"git.sr.ht/~ionous/iffy/parser"
)

// GrammarMaker -
type GrammarMaker interface{ MakeGrammar() GrammarMaker }

func (op *Alias) MakeGrammar() GrammarMaker { return op }

func (op *Directive) MakeGrammar() GrammarMaker { return op }

// acts as AllOf{Words, ... }
func (op *Directive) MakeScanners() (ret parser.Scanner) {
	out := make([]parser.Scanner, len(op.Scans)+1)
	out[0] = parser.Words(op.Lede)
	for i, el := range op.Scans {
		out[i+1] = el.MakeScanner()
	}
	return &parser.AllOf{out}
}
