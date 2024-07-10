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
	Lines [][]match.TokenValue // tokens have pos
	// keeps pronoun context for the most recent line
	// across multiple calls to "schedule"
	// ( if scheduling was channel based, we could put pronounSource on the stack )
	pronouns pronounSource
}

// use the existing tokens as a paragraph
// ( ex. from parsing a plain text section )
func MakeParagraph(file string, lines [][]match.TokenValue) Paragraph {
	return Paragraph{File: file, Lines: lines}
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
func (p *Paragraph) Generate(z weaver.Phase, q Query, u Scheduler) (okay bool, err error) {
	var retry int
	unmatched := p.Lines
	for i, n := range unmatched { // n: is a sentence of tokens
		var best bestMatch
		line := InputState{words: n, pronouns: p.pronouns.nextPronoun()}
		// match a sentence,
		// and if matched Generate/Schedule it for weaving database info
		if matchSentence(z, q, line, &best) {
			lineOfs := n[0].Pos.Y
			source := compact.Source{
				File:    p.File,
				Line:    lineOfs,
				Comment: "a plain-text paragraph",
			}
			// update the paragraph's context so other sentences can refer to it.
			// ( or if no pronoun was matched, or reused, clear it )
			p.pronouns = best.pronouns.nextPronoun()
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
				err = fmt.Errorf("failed to match line %d %s", i, Matched(n).DebugString())
				break
			} else {
				unmatched[retry] = n
				retry++
			}
		}
	}
	// no errors? update the unmatched list
	// fix? a slightly better api would be to return a new paragraph
	// and let the caller check if it has more lines.
	if err == nil {
		if retry > 0 {
			p.Lines = unmatched[:retry]
		} else {
			p.Lines = nil
			okay = true
		}
	}
	return
}
