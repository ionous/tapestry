package jess

import (
	"errors"
	"fmt"

	"git.sr.ht/~ionous/tapestry/lang/compact"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/support/match"
	"git.sr.ht/~ionous/tapestry/weave/weaver"
)

// represents a block of text
type Paragraph struct {
	File string
	// sentences within the paragraph.
	// each sentence is its own slice of tokens.
	// weaving winnows this list.
	Phrases []Phrase
	// index into phrases of unmatched sentences
	unmatched []int
}

// use the existing tokens as a paragraph
// ( ex. from parsing a plain text section )
func MakeParagraph(file string, phrases [][]match.TokenValue) Paragraph {
	return Paragraph{File: file, Phrases: tokensToPhrases(phrases)}
}

// parse the passed string into a paragraph of sentences.
// ( ex. parsing a story DeclareStatement )
func NewParagraph(pos compact.Source, str string, assign rt.Assignment) (ret Paragraph, err error) {
	c := match.Collector{LineOffset: pos.Line, BreakLines: true}
	if e := c.TokenizeString(str); e != nil {
		err = fmt.Errorf("%w reading %s", e, str)
	} else {
		// tack on the assignment
		if cnt := len(c.Tokens); cnt == 0 && assign != nil {
			err = errors.New("unexpected trailing assignment")
		} else if cnt != 0 && assign == nil {
			err = errors.New("expected trailing assignment")
		} else {
			phrases := c.Lines
			if assign != nil && cnt != 0 {
				tv := match.TokenValue{Token: match.Tell, Value: assign}
				c.Tokens = append(c.Tokens, tv)
				phrases = append(phrases, c.Tokens)
			}
			ret = MakeParagraph(pos.File, phrases)
		}
	}
	return
}

// tries to match every (remaining) sentence of the paragraph.
// expects to be called every phase
// returns true when it no longer needs to be called because everything is scheduled
func (p *Paragraph) WeaveParagraph(z weaver.Phase, q Query, u Scheduler) (okay bool, err error) {
	const defaultFlags = 0
	unmatched, retry := p.unmatched, 0
	// first weave initialization
	if unmatched == nil { // ugh. fine for now.
		unmatched = make([]int, len(p.Phrases))
		for i, el := range p.Phrases {
			jc := JessContext{q, u, p, i, defaultFlags}
			in := InputState{words: el.words}
			TryPromisedMatch(jc, in) // launch parallel matches
			unmatched[i] = i
		}
	}
	for _, i := range unmatched {
		jc := JessContext{q, u, p, i, defaultFlags}
		el := &(p.Phrases[i])
		if el.matched != nil {
			continue // parallel matched this.
		}
		var best bestMatch
		line := InputState{words: el.words}
		// match a sentence,
		// and if matched Generate/Schedule it for weaving database info
		if matchSentence(z, jc, line, &best) {
			// track this match
			el.matched = best.match
			if el.topicType == undeterminedTopic {
				el.topicType = otherTopic
			}
			//
			// after matching: try to generate ( which inevitably calls schedule... )
			// ( errors here are critical, and not a request to "retry" )
			if e := best.match.Generate(jc); e != nil {
				err = e
				break
			}
		} else {
			// if it didn't match; retry in a later phase
			// ( but error if we've gone through all the phases without success )
			if z == weaver.NextPhase {
				e := fmt.Errorf("failed to match %s %s %q", p.File, jc.Source().ErrorString(), Matched(el.words).DebugString())
				err = errors.Join(append([]error{e}, el.errs...)...)
				break
			} else {
				unmatched[retry] = i
				retry++
			}
		}
	}
	// no errors? update the unmatched list
	if err == nil {
		p.unmatched = unmatched[:retry]
		okay = retry == 0
	}
	return
}

func TryPromisedMatch(jc JessContext, in InputState) {
	el := jc.CurrentPhrase()
	// property of noun is/are value.
	TryPropertyNounValue(jc, in, el.store, el.reject)
	// noun has property of value.
	TryNounPropertyValue(jc, in, el.store, el.reject)
}
