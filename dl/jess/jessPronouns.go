package jess

// these track author specified names so that plain-english story text
// can refer back to earlier nouns via pronouns
type pronounSource struct {
	// within a paragraph, sentences can establish a source for pronouns.
	// if its not used; its cleared
	source       *Name // for now, singular name; could support plural
	usedPronouns bool  //reset every matching attempt
}

// stored privately in the matched pronoun object
type PronounReference struct {
	source *Name // refers back to whatever was established
}

func (p *pronounSource) isValid() bool {
	return p.usedPronouns
}

// called for every new sentence.
// if a source for pronouns was set, then use that.
func (p *pronounSource) nextPronoun() (ret pronounSource) {
	if p.usedPronouns {
		// ret.usedPronouns will be false
		// until a source is set or a source is referenced.
		ret.source = p.source
	}
	return
}

// at least for now, only works with single nouns
// and the matcher only understands "it"
func (p *pronounSource) setPronounSource(ns Names) {
	if ns.Name != nil && ns.AdditionalNames == nil {
		p.source = ns.Name
		p.usedPronouns = true
	}
}

// called by a specific use of a pronoun ( ex. "it" )
// return true if there was an established name that the pronoun refers to.
// ( and record into the reference what that established name was )
func (p *pronounSource) usePronoun(out *PronounReference) (okay bool) {
	if src := p.source; src != nil {
		p.usedPronouns = true // keep this source alive for another sentence.
		out.source = src
		okay = true
	}
	return
}

// match the *use* of a pronoun ( ex. "it" )
func (op *Pronoun) Match(_ Query, input *InputState) (okay bool) {
	// fix: i think match should be able to return error
	// maybe as a freefunction similar to Optional that takes an error address?
	// or maybe record a status into InputState?
	//
	if width := input.MatchWord(keywords.It); width > 0 && //
		input.pronouns.usePronoun(&op.proref) {
		//
		op.Matched = input.Cut(width)
		*input = input.Skip(width)
		okay = true
	}
	return
}
