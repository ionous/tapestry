package grammar

import (
	"git.sr.ht/~ionous/iffy/dl/composer"
)

var Slots = []composer.Slot{{
	Type: (*GrammarMaker)(nil),
	Desc: "Grammar Maker: Helper for defining parser grammars.",
}, {
	Type: (*ScannerMaker)(nil),
	Desc: "Scanner Maker: Helper for defining input scanners.",
}}

var Slats = []composer.Composer{
	// grammar maker
	(*Alias)(nil),
	(*Directive)(nil),
	// scanner maker
	(*Action)(nil),
	(*AllOf)(nil),
	(*AnyOf)(nil),
	(*Noun)(nil),
	(*Retarget)(nil),
	(*Words)(nil),
	// container for grammar makers
	(*GrammarDecl)(nil),
}
