package grok

type TraitSet struct {
	Kind      Matched
	Traits    []Matched
	WordCount int // including any deleted ands; used during scanning to skip over the complete set of traits
}

func (ts *TraitSet) hasKind() bool {
	return ts.Kind != nil
}

func (ts *TraitSet) kinds() (ret []Matched) {
	if ts.hasKind() {
		ret = []Matched{ts.Kind}
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
