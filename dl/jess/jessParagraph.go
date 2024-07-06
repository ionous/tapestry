package jess

import (
	"errors"
	"fmt"

	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/support/match"
	"git.sr.ht/~ionous/tapestry/weave/mdl"
	"git.sr.ht/~ionous/tapestry/weave/weaver"
)

// represents a block of text
type Paragraph struct {
	file, path string
	lines      [][]match.TokenValue // tokens have pos
}

func (p *Paragraph) NumLines() int {
	return len(p.lines)
}

func MakeParagraph(lines [][]match.TokenValue) Paragraph {
	return Paragraph{lines: lines}
}

func ParagraphPos(pos mdl.Source, str string, assign rt.Assignment) (ret Paragraph, err error) {
	c := match.Collector{BreakLines: true}
	if e := c.Collect(str, pos.Line); e != nil {
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
			ret = Paragraph{
				pos.File, pos.Path, lines,
			}
		}
	}
	return
}

// a function that expects to be called every phase
// returns true when it no longer needs to be called because everything is scheduled
func (p *Paragraph) Generate(z weaver.Phase, q Query, u Scheduler) (okay bool, err error) {
	var retry int
	unmatched := p.lines
	for i, n := range unmatched {
		var best bestMatch
		if matchSentence(z, q, n, &best) {
			// try to generate if matched.
			lineOfs := n[0].Pos.Y
			source := mdl.Source{
				File:    p.file,
				Path:    p.path,
				Line:    lineOfs,
				Comment: "a plain-text paragraph",
			}
			if e := best.match.Generate(Context{q, u, source}); e != nil {
				err = e
				break
			}
		} else {
			// retry if not in final phase
			// otherwise generate an error
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
	if err == nil {
		if retry > 0 {
			p.lines = unmatched[:retry]
		} else {
			p.lines = nil
			okay = true
		}
	}
	return
}
