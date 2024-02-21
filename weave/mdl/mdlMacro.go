package mdl

import "strings"

// the expected amount of names on the left and right hand sides of a macro phrase.
// ( which side is primary and which side is secondary is fairly arbitrary;
//
//	it depends on the definition of the phrase. )
type MacroType int

const (
	// used for "<sources> are kinds of <name of kind>"
	// where the name of the kind is considered a property of the source.
	// ex. `Devices are a kind of prop.`
	Macro_PrimaryOnly MacroType = iota
	// ex. `In the coffin are some coins, a notebook, and the gripping hand.`
	Macro_ManyPrimary
	// ex. `Some coins, a notebook, and the gripping hand are in the coffin.`
	Macro_ManySecondary
	// ex. `Hector and Maria are suspicious of Santa and Santana.`
	Macro_ManyMany
)

type Macro struct {
	Name     string
	Type     MacroType
	Reversed bool
}

func (m MacroType) MultiSource() (okay bool) {
	switch m {
	case Macro_PrimaryOnly, Macro_ManyPrimary, Macro_ManyMany:
		okay = true
	}
	return
}

func (m MacroType) MultiTarget() (okay bool) {
	switch m {
	case Macro_ManySecondary, Macro_ManyMany:
		okay = true
	}
	return
}

// the reversed relation name
func fmtMacro(name string, reversed bool) (ret string) {
	ps := []string{"the"}
	if reversed {
		ps = append(ps, "reversed")
	}
	ps = append(ps, "macro", name)
	return strings.Join(ps, " ")
}
