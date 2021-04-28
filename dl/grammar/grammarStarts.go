package grammar

import (
	"git.sr.ht/~ionous/iffy/dl/composer"
	"git.sr.ht/~ionous/iffy/parser"
)

// GrammarMaker -
type GrammarMaker interface{ MakeGrammar() GrammarMaker }

// Alias allows the user to refer to a noun by one or more other terms.
type Alias struct {
	Names  []string `if:"selector"`
	AsNoun string
}

// Directive starts a parser scanner
type Directive struct {
	Lede  []string `if:"selector"`
	Scans []ScannerMaker
}

// GrammarDecl is the container for all things grammar related.
type GrammarDecl struct {
	Grammar GrammarMaker
}

func (*GrammarDecl) Compose() composer.Spec {
	return composer.Spec{
		Lede:  "grammar",
		Name:  "grammar_decl",
		Spec:  "Understand {grammar:grammar_maker}",
		Slots: []string{"story_statement"},
		Group: "grammar",
		Desc:  `Understand input: Read what the player types and turn it into actions.`,
		Stub:  true,
	}
}

func (*Alias) Compose() composer.Spec {
	return composer.Spec{
		Group:  "grammar",
		Fluent: &composer.Fluid{Role: composer.Function},
	}
}

func (*Directive) Compose() composer.Spec {
	return composer.Spec{
		Group:  "grammar",
		Fluent: &composer.Fluid{Role: composer.Function},
	}
}

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
