package story

import (
	"git.sr.ht/~ionous/tapestry/dl/jess"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"git.sr.ht/~ionous/tapestry/support/jessdb"
	"git.sr.ht/~ionous/tapestry/support/match"
	"git.sr.ht/~ionous/tapestry/weave"
	"git.sr.ht/~ionous/tapestry/weave/weaver"
)

// private member of DeclareStatement
type JessMatches []match.Span

// move into jess/dl
func (op *DeclareStatement) Weave(cat *weave.Catalog) error {
	return cat.Schedule(weaver.LanguagePhase, func(w weaver.Weaves, run rt.Runtime) (err error) {
		if txt, e := safe.GetText(run, op.Text); e != nil {
			err = e
		} else if p, e := jess.NewParagraph(txt.String()); e != nil {
			err = e
		} else {
			q := jessdb.MakeQuery(cat.Modeler, cat.CurrentDomain())
			// a little gross: run a step manually in the language phase
			if ok, e := p.Generate(weaver.LanguagePhase, q, cat); e != nil {
				err = e
			} else if !ok {
				err = cat.Step(func(z weaver.Phase) (bool, error) {
					return p.Generate(z, q, cat)
				})
			}
		}
		return
	})
}
