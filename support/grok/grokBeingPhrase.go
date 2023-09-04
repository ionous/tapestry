package grok

import "github.com/ionous/errutil"

// after finding is/are, check for a macro;
// handle either 1. a pure "is" statement, or 2. an "is-verb" statement.
// 1.  (lhs: The bottle) is (rhs: closed.)
// 2a. (lhs: A device) is (rhs: in the lobby.)
// 2b. (lhs: The coffin) is (rhs: a closed container in the lobby.)
// 2c. (lhs: The coffin) is (rhs: closed in the lobby.)  <-- not my favorite.
//
// tbd: parse the sources when looping over the words ( in the caller? )
func beingPhrase(known Grokker, lhs, rhs []Word) (ret Results, err error) {
	var out Results
	var hack bool
Hack:
	// first, scan for leading traits on the rhs
	// ex. [is] ( rhs: fixed in place .... in the lobby )
	if rightLede, e := parseTraitSet(known, rhs, hack); e != nil {
		err = e
	} else {
		// our nouns.
		sources, targets := &out.Primary, &out.Secondary
		// try to find a macro after the traits:
		afterRightLede := rhs[rightLede.WordCount:]
		if macro, e := known.FindMacro(afterRightLede); e != nil {
			err = e
		} else if len(macro.Name) == 0 {
			// case 1. doesn't have a macro:
			// note: in phrases like "Xs are kinds of Ys"
			// the word "kinds" can get matched to the built in kindsOf.Kind during parseTraitSet
			// even when it should have matched the macro "a kind of"
			// force trait parsing to skip matching any and all kinds, and try again.
			if len(afterRightLede) > 0 {
				if !hack {
					hack = true
					goto Hack
				}
				err = errutil.Fmt("couldnt parse right hand side: %q", Span(afterRightLede).String())
			} else if e := grokNouns(known, &out.Primary, lhs, AllowMany|AllowAnonymous); e != nil {
				err = errutil.New("parsing subjects", e)
			}
		} else {
			// case 2: found a macro:
			out.Macro = macro
			postMacro := afterRightLede[macro.Match.NumWords():]
			var lhsFlag, rhsFlag genFlag
			switch macro.Type {
			case Macro_PrimaryOnly:
				lhsFlag = AllowMany | AllowAnonymous
			case Macro_ManyPrimary:
				lhsFlag = AllowMany | OnlyNamed
				rhsFlag = OnlyOne | AllowAnonymous
			case Macro_ManySecondary:
				lhsFlag = OnlyOne | AllowAnonymous
				rhsFlag = AllowMany | OnlyNamed
			case Macro_ManyMany:
				lhsFlag = AllowMany | OnlyNamed
				rhsFlag = AllowMany | OnlyNamed
			}

			if macro.Reversed {
				sources, targets = targets, sources
				lhsFlag, rhsFlag = rhsFlag, lhsFlag
			}

			// [lhs: The coffin is] (rhs: (pre: a closed container) *in* (post: the antechamber.))
			if e := grokNouns(known, sources, lhs, lhsFlag); e != nil {
				err = errutil.New("parsing subject", e)
			} else {
				// no relation ( kinds of ) don't have secondary noun targets
				// alt: might be to treat the names of kinds in this case as nouns ( and return them in targets )
				if rhsFlag > 0 {
					if e := grokNouns(known, targets, postMacro, rhsFlag); e != nil {
						err = errutil.New("parsing target", e)
					}
				} else {
					// inform specifically denies these right leading traits in this case:
					// [The box is] (right lede: a closed container) kind of (post traits: closed container).
					if rightLede.WordCount > 0 {
						err = makeWordError(rhs[0], "some unexpected kind of properties")
					} else if postMacroTraits, e := ParseTraitSet(known, postMacro); e != nil {
						err = e
					} else if postMacroTraits.WordCount > 0 {
						postMacroTraits.applyTraits(*sources)
					} else if len(postMacro) > 0 {
						if macro.Type != Macro_PrimaryOnly {
							err = errutil.New("unconsumed words", Span(postMacro).String())
						} else {
							// hack part 2 for "are a kind of"
							// should be revisited at some point
							for i, src := range *sources {
								src.Kinds = append(src.Kinds, Span(postMacro))
								(*sources)[i] = src
							}
						}
					}
				}
			}
		}
		if err == nil {
			rightLede.applyTraits(*sources)
		}
	}
	if err == nil {
		ret = out
	}
	return
}
