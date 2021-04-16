package grammar

import (
	"git.sr.ht/~ionous/iffy/dl/composer"
)

var Slots = []composer.Slot{{
	Type: (*ScannerMaker)(nil),
	Desc: "Grammar: Helper for defining parser input scanners.",
}}

var Slats = []composer.Composer{
	(*Action)(nil),
	(*AllOf)(nil),
	(*AnyOf)(nil),
	(*Grammar)(nil),
	(*Noun)(nil),
	(*Retarget)(nil),
	(*Words)(nil),
}

//
type Grammar struct {
	Scanner ScannerMaker
}

func (*Grammar) Compose() composer.Spec {
	return composer.Spec{
		Name:  "grammar_decl",
		Spec:  "Understand {grammar%scanner:scanner_maker}",
		Slots: []string{"story_statement"},
		Group: "grammar",
		Desc:  `Understand grammar: Reading what the player types and turning that text into actions currently is defined with hand written parse trees.`,
		Stub:  true,
	}
}
