package story

import (
	"git.sr.ht/~ionous/tapestry/dl/jess"
	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/lang/compact"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"git.sr.ht/~ionous/tapestry/support/jessdb"
	"git.sr.ht/~ionous/tapestry/weave"
	"git.sr.ht/~ionous/tapestry/weave/weaver"
)

// private member of DeclareStatement
type JessMatches = jess.Paragraph

// DeclareStatements come from both command and plain-text sections of tell files.
// Command sections have *unparsed* text, and naturally decoded assignments.
// Plain-text sections have *parsed* text, and manually decoded assignments.
// Currently, if the command section's text includes a sub-document, i believe jess will panic.
func MakeDeclaration(str string, tail rt.Assignment, ks JessMatches) *DeclareStatement {
	return &DeclareStatement{
		Text:    &literal.TextValue{Value: str},
		Assign:  tail, // tail is here for reference; its already part of the token stream
		matches: ks,
	}
}

// todo: move into jess/dl
func (op *DeclareStatement) Weave(cat *weave.Catalog) error {
	return cat.ScheduleCmd(op, weaver.LanguagePhase,
		func(w weaver.Weaves, run rt.Runtime) (err error) {
			if p, e := op.newParagraph(run); e != nil {
				err = e
			} else {
				// fix: MakeQueryFromPen; shouldn't this be passed in?
				q := jessdb.MakeQuery(cat.Modeler, cat.CurrentScene())
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

// generate a paragraph, or keep whatever was assigned from the flex file.
func (op *DeclareStatement) newParagraph(run rt.Runtime) (ret jess.Paragraph, err error) {
	if m := op.matches; len(m.Lines) > 0 {
		ret = m // already parsed by flexText; keep what's there
	} else if txt, e := safe.GetText(run, op.Text); e != nil {
		err = e
	} else {
		pos := compact.MakeSource(op.GetMarkup(false))
		ret, err = jess.NewParagraph(pos, txt.String(), op.Assign)
	}
	return
}
