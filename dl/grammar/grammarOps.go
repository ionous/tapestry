package grammar

import (
	"git.sr.ht/~ionous/tapestry/parser"
	"git.sr.ht/~ionous/tapestry/support/inflect"
)

// ScannerMaker - creates parser scanners
type ScannerMaker interface{ MakeScanner() parser.Scanner }

func (op *Action) MakeScanner() parser.Scanner {
	return &parser.Action{
		Name: inflect.Normalize(op.Action),
		Args: op.Arguments, // resolved when the action is executed
	}
}

func (op *ChooseOne) MakeScanner() (ret parser.Scanner) {
	if els := op.Options; len(els) == 1 {
		ret = els[0].MakeScanner()
	} else {
		ls := reduce(els)
		ret = &parser.AnyOf{Match: ls}
	}
	return
}

// tbd: i wonder if this should have children ( like any/one of )
// rather than being "inline"
func (op *Focus) MakeScanner() parser.Scanner {
	ls := makeSequence(op.Series)
	return &parser.Focus{Where: op.Player, Match: ls}
}

func (op *Noun) MakeScanner() parser.Scanner {
	var fs parser.Filters
	if k := op.Kind; len(k) > 0 {
		fs = parser.Filters{&parser.HasClass{Name: k}}
	}
	return &parser.Noun{Filters: fs}
}

func (op *Refine) MakeScanner() parser.Scanner {
	ls := reduce(op.Series)
	return &parser.Refine{Match: ls}
}

func (op *Reverse) MakeScanner() parser.Scanner {
	ls := reduce(op.Reverses)
	return &parser.Reverse{Match: ls}
}

func (op *Sequence) MakeScanner() (ret parser.Scanner) {
	return makeSequence(op.Series)
}

func (op *Words) MakeScanner() parser.Scanner {
	// returns an "any of" with individual word matches
	return parser.Words(op.Words)
}

func makeSequence(els []ScannerMaker) (ret parser.Scanner) {
	if len(els) == 1 {
		ret = els[0].MakeScanner()
	} else {
		ret = &parser.AllOf{Match: reduce(els)}
	}
	return
}

func reduce(els []ScannerMaker) []parser.Scanner {
	out := make([]parser.Scanner, len(els))
	for i, el := range els {
		out[i] = el.MakeScanner()
	}
	return out
}
