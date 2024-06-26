package story

import (
	"git.sr.ht/~ionous/tapestry/dl/jess"
	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"git.sr.ht/~ionous/tapestry/support/jessdb"
	"git.sr.ht/~ionous/tapestry/weave"
	"git.sr.ht/~ionous/tapestry/weave/weaver"
)

// private member of DeclareStatement
type JessMatches jess.Paragraph

func MakeDeclaration(str string, tail rt.Assignment, ks JessMatches) *DeclareStatement {
	return &DeclareStatement{
		Text:    &literal.TextValue{Value: str},
		Assign:  tail,
		matches: ks,
	}
}

// todo: move into jess/dl
func (op *DeclareStatement) Weave(cat *weave.Catalog) error {
	return cat.Schedule(weaver.LanguagePhase,
		func(w weaver.Weaves, run rt.Runtime) (err error) {
			if p, e := op.newParagraph(run); e != nil {
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

func (op *DeclareStatement) newParagraph(run rt.Runtime) (ret jess.Paragraph, err error) {
	if m := op.matches; len(m) > 0 {
		ret = jess.Paragraph(m)
	} else if txt, e := safe.GetText(run, op.Text); e != nil {
		err = e
	} else {
		ret, err = jess.NewParagraph(txt.String(), op.Assign)
	}
	return
}
