package grok

type TraitSet struct {
	Kind      Match
	Traits    []Match
	WordCount int // including any deleted ands
}

func (ts *TraitSet) hasKind() bool {
	return ts.Kind != nil
}

func (ts *TraitSet) kinds() (ret []Match) {
	if ts.hasKind() {
		ret = []Match{ts.Kind}
	}
	return ret
}

func (ts *TraitSet) applyTraits(out []Name) {
	if hasKind := ts.hasKind(); hasKind || len(ts.Traits) > 0 {
		for i := range out {
			if hasKind {
				out[i].Kinds = append(out[i].Kinds, ts.Kind)
			}
			out[i].Traits = append(out[i].Traits, ts.Traits...)
		}
	}
}

// public for testing
// ex. "[the] open container"
// ex. "[the] open and openable container"
// ex. "[is] open"
func ParseTraitSet(known Grokker, ws []Word) (out TraitSet, err error) {
	return parseTraitSet(known, ws, false)
}

func parseTraitSet(known Grokker, ws []Word, noKinds bool) (out TraitSet, err error) {
	var scan int
	var prevSep sepFlag
Loop:
	if rest := ws[scan:]; len(rest) > 0 {
		// although its a bit weird englishy-wise
		// inform allows determiners before every trait:
		// ex. The box is an openable and a closed.
		if det, e := known.FindArticle(rest); e != nil {
			err = e
		} else if skipDet := MatchLen(det.Match); skipDet >= len(rest) {
			err = makeWordError(rest[0], "expected some sort of name")
		} else {
			rest = rest[skipDet:]
			if trait, e := known.FindTrait(rest); e != nil {
				err = e
			} else if skipTrait := MatchLen(trait); skipTrait > 0 {
				// eat any ands between traits
				if skipAnd, andSep, e := countAnd(rest[skipTrait:]); e != nil {
					err = e
				} else if skipRest := skipTrait + skipAnd; skipAnd > 0 && skipRest >= len(rest) {
					err = makeWordError(rest[skipTrait], "unexpected trailing separator")
				} else {
					out.Traits = append(out.Traits, trait)
					prevSep = andSep
					scan += skipRest + skipDet
					goto Loop
				}
			} else if prevSep&AndSep == 0 && !noKinds {
				// if it wasn't a trait and some previous trait didnt end with "and",
				// it might be a trailing kind:
				if kind, e := known.FindKind(rest); e != nil {
					err = e
				} else if skipKind := MatchLen(kind); skipKind > 0 {
					out.Kind = kind
					scan += skipKind + skipDet
					// done.
				}
			}
		}
	}
	// after scanning
	if err == nil {
		out.WordCount = scan
	}
	return
}
