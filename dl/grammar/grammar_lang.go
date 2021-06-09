// Code generated by "makeops"; edit at your own risk.
package grammar

import (
	"git.sr.ht/~ionous/iffy/dl/composer"
)

// Action makes a parser scanner producing a script defined action.
type Action struct {
	Action string `if:"label=_,type=string"`
}

var _ ScannerMaker = (*Action)(nil)

func (*Action) Compose() composer.Spec {
	return composer.Spec{
		Name: "action",
		Uses: "flow",
		Lede: "as",
	}
}

// Alias allows the user to refer to a noun by one or more other terms.
type Alias struct {
	Names  []string `if:"label=_,type=string"`
	AsNoun string   `if:"label=as_noun,type=string"`
}

var _ GrammarMaker = (*Alias)(nil)

func (*Alias) Compose() composer.Spec {
	return composer.Spec{
		Name: "alias",
		Uses: "flow",
	}
}

// AllOf makes a parser scanner
type AllOf struct {
	Series []ScannerMaker `if:"label=_"`
}

var _ ScannerMaker = (*AllOf)(nil)

func (*AllOf) Compose() composer.Spec {
	return composer.Spec{
		Name: "all_of",
		Uses: "flow",
	}
}

// AnyOf makes a parser scanner
type AnyOf struct {
	Options []ScannerMaker `if:"label=_"`
}

var _ ScannerMaker = (*AnyOf)(nil)

func (*AnyOf) Compose() composer.Spec {
	return composer.Spec{
		Name: "any_of",
		Uses: "flow",
	}
}

// Directive starts a parser scanner
type Directive struct {
	Lede  []string       `if:"label=_,type=string"`
	Scans []ScannerMaker `if:"label=scans"`
}

var _ GrammarMaker = (*Directive)(nil)

func (*Directive) Compose() composer.Spec {
	return composer.Spec{
		Name: "directive",
		Uses: "flow",
	}
}

// Grammar Read what the player types and turn it into actions.
type Grammar struct {
	Grammar GrammarMaker `if:"label=_"`
}

func (*Grammar) Compose() composer.Spec {
	return composer.Spec{
		Name: "grammar",
		Uses: "flow",
	}
}

// Noun makes a parser scanner
type Noun struct {
	Kind string `if:"label=_,type=string"`
}

var _ ScannerMaker = (*Noun)(nil)

func (*Noun) Compose() composer.Spec {
	return composer.Spec{
		Name: "noun",
		Uses: "flow",
	}
}

// Retarget makes a parser scanner
type Retarget struct {
	Span []ScannerMaker `if:"label=_"`
}

var _ ScannerMaker = (*Retarget)(nil)

func (*Retarget) Compose() composer.Spec {
	return composer.Spec{
		Name: "retarget",
		Uses: "flow",
	}
}

// Reverse makes a parser scanner
type Reverse struct {
	Reverses []ScannerMaker `if:"label=_"`
}

var _ ScannerMaker = (*Reverse)(nil)

func (*Reverse) Compose() composer.Spec {
	return composer.Spec{
		Name: "reverse",
		Uses: "flow",
	}
}

// Self makes a parser scanner which matches the player. ( the player string is just to make the composer happy. )
type Self struct {
	Player string `if:"label=_,type=string"`
}

var _ ScannerMaker = (*Self)(nil)

func (*Self) Compose() composer.Spec {
	return composer.Spec{
		Name: "self",
		Uses: "flow",
	}
}

// Words makes a parser scanner
type Words struct {
	Words []string `if:"label=_,type=string"`
}

var _ ScannerMaker = (*Words)(nil)

func (*Words) Compose() composer.Spec {
	return composer.Spec{
		Name: "words",
		Uses: "flow",
	}
}

var Slots = []interface{}{
	(*GrammarMaker)(nil),
	(*ScannerMaker)(nil),
}
var Slats = []composer.Composer{
	(*Action)(nil),
	(*Alias)(nil),
	(*AllOf)(nil),
	(*AnyOf)(nil),
	(*Directive)(nil),
	(*Grammar)(nil),
	(*Noun)(nil),
	(*Retarget)(nil),
	(*Reverse)(nil),
	(*Self)(nil),
	(*Words)(nil),
}
