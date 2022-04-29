package regen

import "strings"

// represents the str choices, swap picks, or flow terms.
// each spec type uses different members of this struct.
// but the templates only write what they want,
// so a common struct simplifies things.
type Param struct {
	Name, Label, Type, Phrasing string
	Optional, Repeats           bool
}

// true if the label name is required
func (p Param) Labeled() bool {
	return len(p.Label) > 0 &&
		strings.ReplaceAll(p.Label, " ", "_") != p.Name &&
		p.Label != "_" && // anonymous first member, aka an inline flow
		p.Label != "-" // private unserialized data
}

// true if the type name is required
func (p Param) Private() bool {
	return p.Label == "-"
}

func (p Param) Phrased() bool {
	var label string
	if p.Labeled() {
		label = p.Label
	} else {
		label = p.Name
	}
	return len(p.Phrasing) > 0 &&
		strings.ToLower(strings.ReplaceAll(p.Phrasing, " ", "_")) !=
			strings.ToLower(strings.ReplaceAll(label, " ", "_"))
}

// true if the type name is required
func (p Param) Typed() bool {
	return len(p.Type) > 0 && p.Type != p.Name
}

func makeParam(k string, ps map[string]interface{}) Param {
	v := MapOf(k, ps)
	return Param{
		Name:     Detokenize(k),
		Label:    StringOf("tag", v),
		Type:     StringOf("type", v),
		Repeats:  BoolOf("repeats", v),
		Optional: BoolOf("optional", v),
		Phrasing: StringOf("label", v),
	}
}
