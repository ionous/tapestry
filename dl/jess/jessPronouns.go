package jess

import (
	"errors"

	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/weave/weaver"
)

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

// called for every new sentence.
// if a source for pronouns was set, then use that.
func (k *pronounSource) nextPronoun() (ret pronounSource) {
	if k.usedPronouns {
		// ret.usedPronouns will be false
		// until a source is set or a source is referenced.
		ret.source = k.source
	}
	return
}

// at least for now, only works with single nouns
// and the matcher only understands "it"
func (k *pronounSource) setPronounSource(ns Names) {
	if ns.Name != nil && ns.AdditionalNames == nil {
		k.source = ns.Name
		k.usedPronouns = true
	}
}

// called by a specific use of a pronoun ( ex. "it" )
// return true if there was an established name that the pronoun refers to.
// ( and record into the reference what that established name was )
func (k *pronounSource) usePronoun(out *PronounReference) (okay bool) {
	if src := k.source; src != nil {
		k.usedPronouns = true // keep this source alive for another sentence.
		out.source = src
		okay = true
	}
	return
}

// when the query context has "MatchPronouns",
// try to match the *use* of a pronoun ( ex. "it" ).
func (op *Pronoun) Match(q Query, input *InputState) (okay bool) {
	// fix: i think match should be able to return error
	// maybe as a freefunction similar to Optional that takes an error address?
	// or maybe record a status into InputState?
	if matchPronouns(q) {
		if width := input.MatchWord(keywords.It); width > 0 && //
			input.pronouns.usePronoun(&op.proref) {
			//
			op.Matched = input.Cut(width)
			*input = input.Skip(width)
			okay = true
		}
	}
	return
}

// names are often potential nouns;
// this helper generates them as such.
func (op *Pronoun) BuildNouns(q Query, w weaver.Weaves, run rt.Runtime, props NounProperties) (ret []DesiredNoun, err error) {
	if src := op.proref.source; src == nil {
		err = errors.New("missing referenced name")
	} else if n := src.desiredNoun.Noun; len(n) == 0 {
		err = errors.New("missing referenced noun")
	} else {
		// duplicates Noun.BuildNouns:
		if e := writeKinds(w, n, props.Kinds); e != nil {
			err = e
		} else {
			ret = []DesiredNoun{{Noun: n, Traits: props.Traits}}
		}
	}
	return
}
