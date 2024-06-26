package jess

import (
	"errors"

	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/weave/weaver"
)

// fix: maybe move this into weave as a "noun builder" to reduce duplicate queries?
// maybe has an interface -- so that "AddKind" might be immediately adding for mock; but not for real
type DesiredNoun struct {
	Noun    string // key ( once it has been created; or if it already exists. )
	Aliases []string
	Traits  []string
	Values  []DesiredValue
	//
	ArticleTrait      string
	IndefiniteArticle string
}

type NounProperties struct {
	Traits []string
	Kinds  []string
}

type DesiredValue struct {
	Field  string
	Assign rt.Assignment
}

func (n *DesiredNoun) appendTraits(traits []string) {
	n.Traits = append(n.Traits, traits...)
}

func (n *DesiredNoun) appendArticle(a *Article) {
	if a == nil {
		// the lack of a recognized article makes something proper-named.
		// if proper named, it's neither indefinite nor plural
		n.ArticleTrait = ProperNameTrait
	} else {
		if a.flags.Plural {
			n.ArticleTrait = PluralNamedTrait
		}
		if a.flags.Indefinite {
			n.IndefiniteArticle = a.Text
		}
	}
}

// send the contents of the noun to the db
// assumes that the noun key is already set.
func (n DesiredNoun) writeNounValues(w weaver.Weaves) (err error) {
	if e := n.applyAliases(w); e != nil {
		err = e
	} else if e := n.applyArticleTrait(w); e != nil && !errors.Is(e, weaver.Missing) {
		err = e // article traits are considered optional because not all kinds have them.
	} else if e := n.applyTraits(w); e != nil {
		err = e
	} else if e := n.applyArticleValue(w); e != nil && !errors.Is(e, weaver.Missing) {
		err = e
	} else if e := n.applyValues(w); e != nil {
		err = e
	}
	return
}

func (n DesiredNoun) applyArticleTrait(w weaver.Weaves) (err error) {
	if t := n.ArticleTrait; len(t) > 0 {
		err = w.AddNounTrait(n.Noun, t)
	}
	return
}

// the value is split from the trait to match tests
func (n DesiredNoun) applyArticleValue(w weaver.Weaves) (err error) {
	if t := n.IndefiniteArticle; len(t) > 0 {
		err = w.AddNounValue(n.Noun, IndefiniteArticle, text(t, ""))
	}
	return
}

// assumes that the noun key is already set.
func (n DesiredNoun) applyAliases(w weaver.Weaves) (err error) {
	for _, a := range n.Aliases {
		if e := w.AddNounName(n.Noun, a, -1); e != nil {
			err = e
			break
		}
	}
	return
}

// assumes that the noun key is already set.
func (n DesiredNoun) applyTraits(w weaver.Weaves) (err error) {
	for _, t := range n.Traits {
		if e := w.AddNounTrait(n.Noun, t); e != nil {
			err = e
			break
		}
	}
	return
}

// assumes that the noun key is already set.
func (n DesiredNoun) applyValues(w weaver.Weaves) (err error) {
	for _, v := range n.Values {
		if e := w.AddNounValue(n.Noun, v.Field, v.Assign); e != nil {
			err = e
			break
		}
	}
	return
}
