package grok

type traitSet struct {
	kind      []Word
	traits    [][]Word
	wordCount int // including any deleted ands
}

func (ts *traitSet) kinds() (ret [][]Word) {
	if len(ts.kind) > 0 {
		ret = [][]Word{ts.kind}
	}
	return ret
}

func (ts *traitSet) applyTraits(out []Noun) {
	if len(ts.kind) > 0 || len(ts.traits) > 0 {
		for i := range out {
			if len(ts.kind) > 0 {
				out[i].Kinds = append(out[i].Kinds, ts.kind)
			}
			out[i].Traits = append(out[i].Traits, ts.traits...)
		}
	}
}

// ex. "[the] open container"
// ex. "[the] open and openable container"
// ex. "[is] open"
func parseTraitSet(known Grokker, ws []Word) (out traitSet, err error) {
	var scan int
	var prevSep sepFlag
Loop:
	if rest := ws[scan:]; len(rest) > 0 {
		// although its a bit weird englishy-wise
		// inform allows determiners before every trait:
		// ex. The box is an openable and a closed.
		if skipDet := known.FindDeterminer(rest); skipDet >= len(rest) {
			err = makeWordError(rest[0], "expected some sort of name")
		} else {
			rest = rest[skipDet:]
			if skipTrait := known.FindTrait(rest); skipTrait > 0 {
				// eat any ands between traits
				if skipAnd, andSep, e := countAnd(rest[skipTrait:]); e != nil {
					err = e
				} else if skipRest := skipTrait + skipAnd; skipAnd > 0 && skipRest >= len(rest) {
					err = makeWordError(rest[skipTrait], "unexpected trailing separator")
				} else {
					out.traits = append(out.traits, rest[:skipTrait])
					prevSep = andSep
					scan += skipRest + skipDet
					goto Loop
				}
			} else if prevSep&AndSep == 0 {
				// if it wasn't a trait and some previous trait didnt end with "and",
				// it might be a trailing kind:
				if skipKind := known.FindKind(rest); skipKind > 0 {
					out.kind = rest[:skipKind]
					scan += skipKind + skipDet
					// done.
				}
			}
		}
	}
	// after scanning
	if err == nil {
		out.wordCount = scan
	}
	return
}
