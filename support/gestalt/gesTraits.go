package gestalt

import "git.sr.ht/~ionous/tapestry/support/grok"

// matches one or more traits.
// an article can precede every trait;
// a comma or a comma plus an and can follow a trait
// so long as the next span is also a trait.
//
// ex. "[the] open (container...)"
// ex. "[the] open and [an] openable (container...)"
type Traits struct{}

func (*Traits) Match(q Query, cs []InputState) (ret []InputState) {
	// for each input stream
	for _, in := range cs {
		ws := in.Words()
	Loop:
		for i := 0; len(ws) > 0; i++ {
			if and := q.SkipSeparators(ws); and < 0 || (and > 0 && i == 0) {
				break Loop // the first word cant be "and" or ","
			} else if det := q.SkipArticle(ws[and:]); det < 0 {
				break Loop
			} else if trait, width := q.FindTrait(ws[and+det:]); width <= 0 {
				break Loop // trait tested against zero because it must exist
			} else {
				skip := and + det + width
				out := in.Next(skip)
				out.AddResult(MatchedTrait{trait})
				ret = append(ret, out)
				ws = ws[skip:]
			}
		}
	}
	return
}

type MatchedTrait struct {
	grok.Match
}
