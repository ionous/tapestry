package regen

import "strings"

// represents the str choices, swap picks, or flow terms.
// each spec type uses different members of this struct.
// but the templates only write what they want,
// so a common struct simplifies things.
type Param struct {
	Name, Label, Type string
	Optional, Repeats bool
}

// true if the label name is required
func (opt Param) Labeled() bool {
	return len(opt.Label) > 0 &&
		strings.ReplaceAll(opt.Label, " ", "_") != opt.Name &&
		opt.Label != "_" && // anonymous first member, aka a trim flow
		opt.Label != "-" // private unserialized data
}

// true if the type name is required
func (opt Param) Private() bool {
	return opt.Label == "-"
}

// true if the type name is required
func (opt Param) Typed() bool {
	return len(opt.Type) > 0 && opt.Type != opt.Name
}
