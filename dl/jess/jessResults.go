package jess

import (
	"git.sr.ht/~ionous/tapestry/rt"
)

// fix: maybe move this into weave as a "noun builder" to reduce duplicate queries?
// maybe has an interface -- so that "AddKind" might be immediately adding for mock; but not for real
type DesiredNoun struct {
	Noun    string // key ( once it has been created; or if it already exists. )
	Aliases []string
	Traits  []string
	Values  []DesiredValue
}

type DesiredValue struct {
	Field  string
	Assign rt.Assignment
}

func (n *DesiredNoun) addValue(field string, assign rt.Assignment) {
	n.Values = append(n.Values, DesiredValue{field, assign})
}

func (n *DesiredNoun) addArticle(a *Article) {
	if a == nil {
		// the lack of a recognized article makes something proper-named.
		n.Traits = append([]string{ProperNameTrait}, n.Traits...)
	} else {
		if a.Flags.Plural {
			n.Traits = append([]string{PluralNamedTrait}, n.Traits...)
		}
		if a.Flags.Indefinite {
			n.addValue(IndefiniteArticle, text(a.Text, ""))
		}
	}
}

// send the contents of the noun to the db
// assumes that the noun key is already set.
func (n DesiredNoun) generateValues(rar *Context) (err error) {
	if e := n.applyAliases(rar); e != nil {
		err = e
	} else if e := n.applyTraits(rar); e != nil {
		err = e
	} else if e := n.applyValues(rar); e != nil {
		err = e
	}
	return
}

// assumes that the noun key is already set.
func (n DesiredNoun) applyAliases(rar *Context) (err error) {
	for _, a := range n.Aliases {
		if e := rar.AddNounName(n.Noun, a, -1); e != nil {
			err = e
			break
		}
	}
	return
}

// assumes that the noun key is already set.
func (n DesiredNoun) applyTraits(rar *Context) (err error) {
	for _, t := range n.Traits {
		if e := rar.AddNounTrait(n.Noun, t); e != nil {
			err = e
			break
		}
	}
	return
}

// assumes that the noun key is already set.
func (n DesiredNoun) applyValues(rar *Context) (err error) {
	for _, v := range n.Values {
		if e := rar.AddNounValue(n.Noun, v.Field, v.Assign); e != nil {
			err = e
			break
		}
	}
	return
}
