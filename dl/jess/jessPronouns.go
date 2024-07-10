package jess

// these track "author facing" pronouns so that plain-english story text
// can refer to earlier nouns via pronouns
type currentPronoun struct {
	// within a paragraph, sentences can establish a source for pronouns.
	// if its not used; its cleared
	source       *Name // for now, singular name; could support plural
	usedPronouns bool  //reset every matching attempt
}

type pronounSource struct {
	source       *Name
	usedPronouns bool
}

// stored privately in the matched pronoun object
type PronounReference struct {
	source *Name // refers back to whatever was established
}

// called by a specific use of a pronoun ( ex. "it" )
// return true if there was an established name that the pronoun might refer to.
// ( and record into the reference what that established name was )
func (p *currentPronoun) usePronoun(out *PronounReference) (okay bool) {
	if src := p.source; src != nil {
		out.source = src
		okay = true
	}
	return
}

func (p *currentPronoun) nextPronoun() (ret pronounSource) {
	if p.usedPronouns {
		ret.source = p.source
		p.usedPronouns = false // was used, then clear for next time.
	} else {
		(*p) = currentPronoun{} // not used, clear entirely.
	}
	return
}

// at least for now, only works with single nouns
// and the matcher only understands "it"
func (p *pronounSource) setPronous(ns Names) {
	if ns.Name != nil && ns.AdditionalNames == nil {
		p.source = ns.Name
		p.usedPronouns = true
	}
}

// fix: i think match should be able to return error
// maybe as a freefunction similar to Optional that takes an error address?
// or maybe record a status into paragrapH?
// there is probably something wrong here....
// if a parent matched a pronoun source --
// we should copy that data into it.
// so that it can be unwound if some other branch doesnt
func (op *Pronoun) Match(_ Query, input *InputState) (okay bool) {
	// if width := input.MatchWord(keywords.It); width > 0 && //
	// 	//p.currentPronoun.usePronoun(&op.proref) {
	// 	//
	// 	op.Matched = input.Cut(width)
	// 	*input = input.Skip(width)
	// 	okay = true
	// }
	return
}
