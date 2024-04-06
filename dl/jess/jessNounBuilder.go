package jess

import (
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/weave/weaver"
)

// really this builds "pending nouns"....
// to support "counted nouns" any given specification can generate multiple nouns
// ( even though all, other than "names" and "counted nouns" only generate one a piece. )
type NounBuilder interface {
	BuildNouns(Query, weaver.Weaves, rt.Runtime, NounProperties) ([]DesiredNoun, error)
}

// useful for dispatching a parent's call to build nouns to one of its matched children.
func buildNounsFrom(q Query, w weaver.Weaves, run rt.Runtime, props NounProperties, options ...nounBuilderRef) (ret []DesiredNoun, err error) {
	for _, opt := range options {
		if !opt.IsNil {
			ret, err = opt.BuildNouns(q, w, run, props)
			break
		}
	}
	return
}

func buildAnon(w weaver.Weaves, plural, singular string, props NounProperties) (ret DesiredNoun, err error) {
	n := w.GenerateUniqueName(singular)
	if e := w.AddNounKind(n, plural); e != nil {
		err = e // all errors, including duplicates would be bad here.
	} else if e := w.AddNounName(n, n, 0); e != nil {
		err = e // ^ so authors can refer to it by the dashed name
	} else if e := writeKinds(w, n, props.Kinds); e != nil {
		err = e // any *additional* kinds.
	} else {
		ret = DesiredNoun{
			// no name and no article because, the object itself is anonymous.
			// ( the article associated with the kind gets eaten )
			Noun:    n,
			Aliases: []string{singular}, // at runtime, "triangle" means "triangle-1"
			Traits:  append([]string{CountedTrait}, props.Traits...),
			Values: []DesiredValue{{
				// to print "triangle-1" as "triangle"
				PrintedName, text(singular, ""),
			}},
		}
	}
	return
}
