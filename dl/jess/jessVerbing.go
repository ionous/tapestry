package jess

// -------------------------------------------------------------------------
// VerbNamesAreNames
// -------------------------------------------------------------------------

func (op *VerbNamesAreNames) GetNouns() Names {
	return op.Names
}
func (op *VerbNamesAreNames) GetOtherNouns() Names {
	return op.OtherNames
}
func (op *VerbNamesAreNames) GetAdjectives() (_ Adjectives) {
	return
}
func (op *VerbNamesAreNames) GetMacro() (ret Macro) {
	return op.Verb.Macro
}
func (op *VerbNamesAreNames) IsReversed() bool {
	return !op.Verb.Macro.Reversed
}
func (op *VerbNamesAreNames) Generate(rar *Context) error {
	return generateVerbPhrase(rar, op)
}
func (op *VerbNamesAreNames) Match(q Query, input *InputState) (okay bool) {
	if next := *input; //
	op.Verb.Match(q, &next) &&
		op.Names.Match(q, &next) &&
		op.Are.Match(q, &next) &&
		op.OtherNames.Match(q, &next) {
		*input, okay = next, true
	}
	return
}

// -------------------------------------------------------------------------
// NamesVerbNames
// -------------------------------------------------------------------------

func (op *NamesVerbNames) GetNouns() Names {
	return op.Names
}
func (op *NamesVerbNames) GetOtherNouns() Names {
	return op.OtherNames
}
func (op *NamesVerbNames) GetAdjectives() (_ Adjectives) {
	return
}
func (op *NamesVerbNames) GetMacro() (ret Macro) {
	return op.Verb.Macro
}
func (op *NamesVerbNames) IsReversed() bool {
	return op.Verb.Macro.Reversed
}
func (op *NamesVerbNames) Generate(rar *Context) error {
	return generateVerbPhrase(rar, op)
}
func (op *NamesVerbNames) Match(q Query, input *InputState) (okay bool) {
	if next := *input; //
	op.Names.Match(q, &next) &&
		!op.Names.HasAnonymousKind() &&
		op.Are.Match(q, &next) &&
		op.Verb.Match(q, &next) &&
		op.OtherNames.Match(q, &next) {
		*input, okay = next, true
	}
	return
}

// -------------------------------------------------------------------------
// NamesAreLikeVerbs
// -------------------------------------------------------------------------

func (op *NamesAreLikeVerbs) GetNouns() Names {
	return op.Names
}
func (op *NamesAreLikeVerbs) GetOtherNouns() (ret Names) {
	if v := op.VerbPhrase; v != nil {
		ret = v.PlainNames
	}
	return
}
func (op *NamesAreLikeVerbs) GetAdjectives() Adjectives {
	return op.Adjectives
}
func (op *NamesAreLikeVerbs) GetMacro() (ret Macro) {
	if v := op.VerbPhrase; v != nil {
		ret = v.Verb.Macro
	}
	return
}
func (op *NamesAreLikeVerbs) IsReversed() (okay bool) {
	if v := op.VerbPhrase; v != nil {
		okay = v.Verb.Macro.Reversed
	}
	return
}
func (op *NamesAreLikeVerbs) Generate(rar *Context) error {
	return generateVerbPhrase(rar, op)
}
func (op *NamesAreLikeVerbs) Match(q Query, input *InputState) (okay bool) {
	if next := *input; //
	op.Names.Match(q, &next) &&
		!op.Names.HasAnonymousKind() &&
		op.Are.Match(q, &next) &&
		op.Adjectives.Match(q, &next) {
		Optional(q, &next, &op.VerbPhrase)
		*input, okay = next, true
	}
	return
}

// -------------------------------------------------------------------------
// jessVerbPhrase
// helper to generate nouns;
// could be a struct, or even multiple parameters instead of a pull/getter
// not sure what's best: can other phrases use the same patter... whatever it is?
// ( esp. seeing as they all need to schedule to the same phase )
// -------------------------------------------------------------------------

// fix? this interface means that Names can contain zero matches.
type jessVerbPhrase interface {
	GetNouns() Names
	GetOtherNouns() Names
	GetAdjectives() Adjectives
	GetMacro() Macro
	IsReversed() bool
}

func generateVerbPhrase(rar *Context, p jessVerbPhrase) error {
	return rar.PostProcess(GenerateNouns, func(q Query) (err error) {
		if ts, ks, e := p.GetAdjectives().Reduce(); e != nil {
			err = e
		} else if lhs, e := p.GetNouns().BuildNouns(q, rar, ts, ks); e != nil {
			err = e
		} else if rhs, e := p.GetOtherNouns().BuildNouns(q, rar, nil, nil); e != nil {
			err = e
		} else {
			macro := p.GetMacro()
			if p.IsReversed() {
				lhs, rhs = rhs, lhs
			}
			err = genNouns(rar, lhs, rhs, func(src, tgt []DesiredNoun) (err error) {
				// note: some phrases "the box is open" dont have macros.
				// in that case, genNouns itself does all the work.
				if len(macro.Name) > 0 {
					err = rar.Apply(macro, reduceNouns(src), reduceNouns(tgt))
				}
				return
			})
		}
		return
	})
}

// fix: this seems silly
func reduceNouns(ns []DesiredNoun) []string {
	out := make([]string, len(ns))
	for i, el := range ns {
		out[i] = el.Noun
	}
	return out
}
