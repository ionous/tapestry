package grammar

import (
	"git.sr.ht/~ionous/tapestry/parser"
)

// ScannerMaker - creates parser scanners
type ScannerMaker interface{ MakeScanner() parser.Scanner }

func (op *Action) MakeScanner() parser.Scanner {
	return &parser.Action{op.Action}
}

func (op *AllOf) MakeScanner() (ret parser.Scanner) {
	if els := op.Series; len(els) == 1 {
		ret = els[0].MakeScanner()
	} else {
		ls := reduce(els)
		ret = &parser.AllOf{ls}
	}
	return
}

func (op *AnyOf) MakeScanner() (ret parser.Scanner) {
	if els := op.Options; len(els) == 1 {
		ret = els[0].MakeScanner()
	} else {
		ls := reduce(els)
		ret = &parser.AnyOf{ls}
	}
	return
}

func (op *Noun) MakeScanner() parser.Scanner {
	var fs parser.Filters
	if k := op.Kind; len(k) > 0 {
		fs = parser.Filters{&parser.HasClass{k}}
	}
	return &parser.Noun{fs}
}

func (op *Retarget) MakeScanner() parser.Scanner {
	ls := reduce(op.Span)
	return &parser.Target{ls}
}

func (op *Reverse) MakeScanner() parser.Scanner {
	ls := reduce(op.Reverses)
	return &parser.Reverse{ls}
}

func (op *Self) MakeScanner() parser.Scanner {
	return &parser.Eat{&parser.Focus{"self", &parser.Noun{}}}
}

func (op *Words) MakeScanner() parser.Scanner {
	// returns an "any of" with individual word matches
	return parser.Words(op.Words)
}

func reduce(els []ScannerMaker) []parser.Scanner {
	out := make([]parser.Scanner, len(els))
	for i, el := range els {
		out[i] = el.MakeScanner()
	}
	return out
}
