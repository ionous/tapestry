package story

import (
	"git.sr.ht/~ionous/tapestry/dl/jess"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"git.sr.ht/~ionous/tapestry/support/jessdb"
	"git.sr.ht/~ionous/tapestry/weave"
	"git.sr.ht/~ionous/tapestry/weave/weaver"
)

// Execute - called by the macro runtime during weave.
func (op *DeclareStatement) Execute(macro rt.Runtime) error {
	return Weave(macro, op)
}

func (op *DeclareStatement) Weave(cat *weave.Catalog) error {
	return cat.Schedule(weaver.LanguagePhase, func(w weaver.Weaves, run rt.Runtime) (err error) {
		if txt, e := safe.GetText(run, op.Text); e != nil {
			err = e
		} else if p, e := jess.NewParagraph(txt.String()); e != nil {
			err = e
		} else {
			q := jessdb.MakeQuery(cat.Modeler, cat.CurrentDomain())
			err = cat.Step(func(z weaver.Phase) (bool, error) {
				return p.Generate(z, q, cat)
			})
		}
		return
	})
}
