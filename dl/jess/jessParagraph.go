package jess

import (
	"errors"
	"fmt"

	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/support/match"
	"git.sr.ht/~ionous/tapestry/weave/weaver"
)

// represents a block of text
type Paragraph [][]match.TokenValue

func NewParagraph(str string, assign rt.Assignment) (ret Paragraph, err error) {
	c := match.Collector{BreakLines: true}
	if e := c.Collect(str); e != nil {
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
			ret = lines
		}
	}
	return
}

// a function that expects to be called every phase
// returns true when it no longer needs to be called because everything is scheduled
func (p *Paragraph) Generate(z weaver.Phase, q Query, u Scheduler) (okay bool, err error) {
	var retry int
	unmatched := (*p)
	for i, n := range unmatched {
		var best bestMatch
		if matchSentence(z, q, n, &best) {
			// try to generate if matched.
			if e := best.match.Generate(Context{q, u}); e != nil {
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
			(*p) = unmatched[:retry]
		} else {
			(*p) = nil
			okay = true
		}
	}
	return
}
