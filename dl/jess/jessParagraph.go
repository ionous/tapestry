package jess

import (
	"fmt"

	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/support/match"
	"git.sr.ht/~ionous/tapestry/weave/weaver"
)

// represents a block of text;
// holds individual sentences of the paragraph.
type Paragraph struct {
	Spans  []match.Span
	Assign rt.Assignment
}

func NewParagraph(str string, assign rt.Assignment) (ret Paragraph, err error) {
	if spans, e := match.MakeSpans(str); e != nil {
		err = fmt.Errorf("%w reading %s", e, str)
	} else {
		ret = Paragraph{
			Spans:  spans,
			Assign: assign,
		}
	}
	return
}

// a function that expects to be called every phase
// returns true when it no longer needs to be called because everything is scheduled
func (p *Paragraph) Generate(z weaver.Phase, q Query, u Scheduler) (okay bool, err error) {
	var retry int
	unmatched := p.Spans
	for i, n := range unmatched {
		var best bestMatch
		next := MakeInput(n, p.Assign)
		if matchSentence(z, q, next, &best) {
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
				unmatched[retry] = n
				retry++
			}
		}
	}
	// no errors? update the unmatched list
	if err == nil {
		if retry > 0 {
			p.Spans = unmatched[:retry]
		} else {
			p.Spans = nil
			okay = true
		}
	}
	return
}
