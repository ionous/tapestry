package grammar

import (
	"git.sr.ht/~ionous/iffy/dl/composer"
	"git.sr.ht/~ionous/iffy/parser"
)

// Scanner - mirrors package parser for composing grammars
type ScannerMaker interface{ MakeScanner() parser.Scanner }

type Action struct {
	Action string `if:"selector"`
}

type AllOf struct {
	Series []ScannerMaker `if:"selector"`
}

type AnyOf struct {
	Options []ScannerMaker `if:"selector"`
}

type Noun struct {
	Kind string `if:"selector"`
}

type Retarget struct {
	Span []ScannerMaker `if:"selector"`
}

type Words struct {
	Words string `if:"selector"`
}

func (*Action) Compose() composer.Spec {
	return composer.Spec{
		Fluent: &composer.Fluid{Name: "as", Role: composer.Function},
		Group:  "grammar",
	}
}

func (*AllOf) Compose() composer.Spec {
	return composer.Spec{
		Group:  "grammar",
		Fluent: &composer.Fluid{Role: composer.Function},
	}
}

func (*AnyOf) Compose() composer.Spec {
	return composer.Spec{
		Group:  "grammar",
		Fluent: &composer.Fluid{Role: composer.Function},
	}
}

func (*Noun) Compose() composer.Spec {
	return composer.Spec{
		Group:  "grammar",
		Fluent: &composer.Fluid{Role: composer.Function},
	}
}

func (*Retarget) Compose() composer.Spec {
	return composer.Spec{
		Group:  "grammar",
		Fluent: &composer.Fluid{Role: composer.Function},
	}
}

func (*Words) Compose() composer.Spec {
	return composer.Spec{
		Group:  "grammar",
		Fluent: &composer.Fluid{Role: composer.Function},
	}
}

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
		ret = &parser.AllOf{ls}
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

func (op *Words) MakeScanner() parser.Scanner {
	return parser.Words(op.Words)
}

func reduce(els []ScannerMaker) []parser.Scanner {
	out := make([]parser.Scanner, len(els))
	for i, el := range els {
		out[i] = el.MakeScanner()
	}
	return out
}