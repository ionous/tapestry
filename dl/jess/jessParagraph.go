package jess

import (
	"fmt"

	"git.sr.ht/~ionous/tapestry/support/match"
	"git.sr.ht/~ionous/tapestry/weave/mdl"
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

// returns true when completely consumed;
// caller no longer needs to process subsequent phrases
func (p *Paragraph) Generate(ctx *Context, z Phase) (okay bool, err error) {
	var retry int
	for _, u := range p.unmatched {
		var best bestMatch
		if matchSentence(ctx, z, u, &best) {
			// try to generate if matched.
			if e := best.match.Generate(ctx); e != nil {
				err = e
				break
			}
		} else {
			// retry if not in final phase
			// otherwise generate an error
			if z == mdl.FinalPhase {
				err = fmt.Errorf("failed to match %s", u.String())
				break
			} else {
				p.unmatched[retry] = u
				retry++
			}
		}
	}
	// no errors? update the unmatched list
	if err == nil {
		p.unmatched = p.unmatched[:retry]
		okay = true
	}
	return
}
