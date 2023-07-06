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
	// first, scan for leading traits on the rhs
	// ex. [is] ( rhs: fixed in place .... in the lobby )
	if rightLede, e := ParseTraitSet(known, rhs); e != nil {
		err = e
	} else {
		sources, targets := &out.Sources, &out.Targets

		// try to find a macro after the traits:
		afterRightLede := rhs[rightLede.WordCount:]
		if macro, e := known.FindMacro(afterRightLede); e != nil {
			err = e
		} else if len(macro.Name) == 0 {
			// case 1. doesn't have a macro:
			if e := grokNouns(known, &out.Sources, lhs, AllowMany|AllowAnonymous); e != nil {
				err = errutil.New("parsing subjects", e)
			}
		} else {
			// case 2: found a macro:
			out.Macro = macro
			postMacro := afterRightLede[macro.Match.NumWords():]
			var lhsFlag, rhsFlag genFlag
			switch macro.Type {
			case Macro_SourcesOnly:
				lhsFlag = AllowMany | OnlyNamed
			case Macro_ManySources:
				lhsFlag = AllowMany | OnlyNamed
				rhsFlag = OnlyOne | AllowAnonymous
			case Macro_ManyTargets:
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
					} else {
						postMacroTraits.applyTraits(*sources)
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
