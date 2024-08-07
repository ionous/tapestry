package jess

import (
	"errors"
	"fmt"
	"log"

	"git.sr.ht/~ionous/tapestry/lang/compact"
	"git.sr.ht/~ionous/tapestry/lang/typeinfo"
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
	Lines []Line
	// keeps pronoun context for the most recent line
	// across multiple calls to "schedule"
	// ( if scheduling was channel based, we could put pronounSource on the stack )
	pronouns pronounSource
	//
	unmatched []*Line
}

type Line struct {
	words []match.TokenValue // fix? all tokens have pos; we only really need the first.
	// the successful match; mainly for debugging; its already written itself to the database
	matched typeinfo.Instance // could store MatchedPhrase maybe.
	//
	errs []error // things we've tried
}

// use the existing tokens as a paragraph
// ( ex. from parsing a plain text section )
func MakeParagraph(file string, lines [][]match.TokenValue) Paragraph {
	return Paragraph{File: file, Lines: linesToLines(lines)}
}

// fix: make collector speak in terms of lines?
func linesToLines(src [][]match.TokenValue) []Line {
	out := make([]Line, len(src))
	for i, el := range src {
		if len(el) == 0 {
			panic("?")
		}
		out[i].words = el
	}
	return out
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
			lines := c.Lines
			if assign != nil && cnt != 0 {
				tv := match.TokenValue{Token: match.Tell, Value: assign}
				c.Tokens = append(c.Tokens, tv)
				lines = append(lines, c.Tokens)
			}
			ret = MakeParagraph(pos.File, lines)
		}
	}
	return
}

// tries to match every (remaining) sentence of the paragraph.
// expects to be called every phase
// returns true when it no longer needs to be called because everything is scheduled
func (p *Paragraph) WeaveParagraph(z weaver.Phase, q Query, u Scheduler) (okay bool, err error) {
	unmatched, retry := p.unmatched, 0
	// first weave initialization
	if unmatched == nil { // ugh. fine for now.
		unmatched = make([]*Line, len(p.Lines))
		for i := range p.Lines {
			el := &(p.Lines[i])
			unmatched[i] = el              // fine so long as we don't resize
			el.StartParallelMatch(p, q, u) // launch parallel matches
		}
	}
	for i, el := range unmatched {
		_ = i // i is useful for debugging.
		if el.matched != nil {
			continue // parallel matched this.
		}
		var best bestMatch
		line := InputState{p: p, words: el.words, pronouns: p.pronouns.nextPronoun()}
		// match a sentence,
		// and if matched Generate/Schedule it for weaving database info
		if matchSentence(z, q, line, &best) {
			// track this match
			el.matched = best.match
			//
			source := line.Source()
			// update the paragraph's context so other sentences can refer to it.
			// ( or if no pronoun was matched, or reused, clear it )
			p.pronouns = best.pronouns
			// after matching: try to generate ( which inevitably calls schedule... )
			// ( errors here are critical, and not a request to "retry" )
			if e := best.match.Generate(Context{q, u, source}); e != nil {
				err = e
				break
			}
		} else {
			// if it didn't match; retry in a later phase
			// ( but error if we've gone through all the phases without success )
			if z == weaver.NextPhase {
				e := fmt.Errorf("failed to match %s %s %q", p.File, line.Source().ErrorString(), Matched(el.words).DebugString())
				err = errors.Join(append([]error{e}, el.errs...)...)
				break
			} else {
				unmatched[retry] = el
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

func (el *Line) StartParallelMatch(p *Paragraph, q Query, u Scheduler) {
	jc := JessContext{q, u}
	in := InputState{p: p, words: el.words}
	// property of noun is/are value.
	TryPropertyNounValue(jc, in, el.store, el.reject)
	// noun has property of value.
	TryNounPropertyValue(jc, in, el.store, el.reject)
}

// fix? maybe there's a difference b/t FailedMatch and other errors?
// Failed indicates the "shape" is wrong
// other errors indicates the content of the particular shape is wrong.
func (el *Line) store(res ParallelMatcher) {
	if el.matched != nil {
		log.Println("matched multiple phases")
	} else {
		el.matched = res
	}
}

func (el *Line) reject(e error) {
	el.errs = append(el.errs, e)
}
