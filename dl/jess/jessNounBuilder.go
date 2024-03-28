package jess

// really this builds "pending nouns"....
// to support "counted nouns" any given specification can generate multiple nouns
// ( even though all, other than "names" and "counted nouns" only generate one a piece. )
type NounBuilder interface {
	BuildNouns(*Context, NounProperties) ([]DesiredNoun, error)
}

// useful for dispatching a parent's call to build nouns to one of its matched children.
func buildNounsFrom(ctx *Context, props NounProperties, options ...nounBuilderRef) (ret []DesiredNoun, err error) {
	for _, opt := range options {
		if !opt.IsNil {
			ret, err = opt.BuildNouns(ctx, props)
			break
		}
	}
	return
}

func buildAnon(rar *Context, plural, singular string, props NounProperties) (ret DesiredNoun, err error) {
	n := rar.GenerateUniqueName(singular)
	if e := rar.AddNounName(n, n, 0); e != nil {
		err = e // ^ so authors can refer to it by the dashed name
	} else if e := rar.AddNounKind(n, plural); e != nil {
		err = e // all errors, including duplicates would be bad here.
	} else if e := writeKinds(rar, n, props.Kinds); e != nil {
		err = e // any *additional* kinds.
	} else {
		ret = DesiredNoun{
			// no name and no article because, the object itself is anonymous.
			// ( the article associated with the kind gets eaten )
			Noun:    n,
			Aliases: []string{singular}, // at runtime, "triangle" means "triangles-1"
			Traits:  append([]string{CountedTrait}, props.Traits...),
			Values: []DesiredValue{{
				// to print "triangles-1" as "triangle"
				PrintedName, text(singular, ""),
			}},
		}
	}
	return
}
