package jess

// really this builds "pending nouns"....
// to support "counted nouns" any given specification can generate multiple nouns
// ( even though all, other than "names" and "counted nouns" only generate one a piece. )
type NounBuilder interface {
	BuildNouns(q Query, rar Registrar, traits, kinds []string) ([]DesiredNoun, error)
}

// useful for dispatching a parent's call to build nouns to one of its matched children.
func buildNounsFrom(q Query, rar Registrar, ts, ks []string, options ...nounBuilderRef) (ret []DesiredNoun, err error) {
	for _, opt := range options {
		if !opt.IsNil {
			ret, err = opt.BuildNouns(q, rar, ts, ks)
			break
		}
	}
	return
}

func buildAnon(rar Registrar, plural, singular string, ts, ks []string) (ret DesiredNoun, err error) {
	n := rar.GetUniqueName(singular)
	if e := rar.AddNounKind(n, plural); e != nil {
		err = e // all errors, including duplicates would be bad here.
	} else if e := rar.AddNounName(n, n, 0); e != nil {
		err = e // ^ so authors can refer to it by the dashed name
	} else if e := registerKinds(rar, n, ks); e != nil {
		err = e // any *additional* kinds.
	} else {
		ret = DesiredNoun{
			// no name and no article because, the object itself is anonymous.
			// ( the article associated with the kind gets eaten )
			Noun:    n,
			Aliases: []string{singular}, // at runtime, "triangle" means "triangles-1"
			Traits:  append([]string{CountedTrait}, ts...),
			Values: []DesiredValue{{
				// to print "triangles-1" as "triangle"
				PrintedName, text(singular, ""),
			}},
		}
	}
	return
}
