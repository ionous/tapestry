package jess

import (
	"fmt"

	"git.sr.ht/~ionous/tapestry/support/match"
	"git.sr.ht/~ionous/tapestry/weave/weaver"
)

// represents a block of text
type Paragraph struct {
	// individual sentences of the paragraph.
	unmatched []match.Span
}

func NewParagraph(str string) (ret *Paragraph, err error) {
	if spans, e := match.MakeSpans(str); e != nil {
		err = fmt.Errorf("%w reading %s", e, str)
	} else {
		ret = &Paragraph{unmatched: spans}
	}
	return
}

// a function that expects to be called every phase
// returns true when it no longer needs to be called because everything is scheduled
func (p *Paragraph) Generate(z weaver.Phase, q Query, u Scheduler) (okay bool, err error) {
	var retry int
	for i, n := range p.unmatched {
		var best bestMatch
		if matchSentence(Context{q, u}, z, n, &best) {
			// try to generate if matched.
			if e := best.match.Generate(Context{q, u}); e != nil {
				err = e
				break
			}
		} else {
			// retry if not in final phase
			// otherwise generate an error
			if z == weaver.NextPhase {
				err = fmt.Errorf("failed to match line %d %s", i, n.String())
				break
			} else {
				p.unmatched[retry] = n
				retry++
			}
		}
	}
	// no errors? update the unmatched list
	if err == nil {
		if retry > 0 {
			p.unmatched = p.unmatched[:retry]
		} else {
			p.unmatched = nil
			okay = true
		}
	}
	return
}
