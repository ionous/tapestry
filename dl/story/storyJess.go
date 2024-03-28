package story

import (
	"errors"

	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/dl/jess"
	"git.sr.ht/~ionous/tapestry/rt"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"git.sr.ht/~ionous/tapestry/support/jessdb"
	"git.sr.ht/~ionous/tapestry/weave"
	"git.sr.ht/~ionous/tapestry/weave/mdl"
	"github.com/ionous/errutil"
)

// Execute - called by the macro runtime during weave.
func (op *DeclareStatement) Execute(macro rt.Runtime) error {
	return Weave(macro, op)
}

func (op *DeclareStatement) Weave(cat *weave.Catalog) error {
	return cat.Schedule(weave.LanguagePhase, func(w *weave.Weaver) (err error) {
		if txt, e := safe.GetText(w, op.Text); e != nil {
			err = e
		} else if p, e := jess.NewParagraph(txt.String()); e != nil {
			err = e
		} else {
			// fix: one context per domain?
			q := jessdb.MakeQuery(cat.Modeler, w.Domain)
			ja := jessAdapter{w, w.Pin()}
			ctx := jess.NewContext(q, ja)
			err = cat.Step(func(z weave.Phase) (err error) {
				if _, e := p.Generate(ctx, z); e != nil {
					err = e
				} else if e := ctx.UpdatePhase(z); e != nil {
					err = e
				}
				return
			})
		}
		return
	})
}

// fix... obviously.
// mostly its the same as Pen, but there are some overrides, and rename *sigh*
type jessAdapter struct {
	w *weave.Weaver
	*mdl.Pen
}

func (ja jessAdapter) AddNounTrait(noun, trait string) error {
	return ja.w.AddNounValue(ja.Pen, noun, trait, truly())
}

// fix: shouldn't there also be an AddNounPath? ( better to merge the two if at all possible )
func (ja jessAdapter) AddNounValue(noun, prop string, val rt.Assignment) error {
	return ja.w.AddNounValue(ja.Pen, noun, prop, val)
}
func (ja jessAdapter) GetClosestNoun(name string) (string, error) {
	return ja.w.GetClosestNoun(name)
}
func (ja jessAdapter) GetPlural(word string) string {
	return ja.w.PluralOf(word)
}
func (ja jessAdapter) GetSingular(word string) string {
	return ja.w.SingularOf(word)
}
func (ja jessAdapter) GenerateUniqueName(category string) string {
	return ja.w.Catalog.NewCounter(category)
}
func (ja jessAdapter) AddFact(key string, parts ...string) (err error) {
	if ok, e := ja.Pen.AddFact(key, parts...); e != nil {
		err = e
	} else if !ok {
		err = mdl.Duplicate
	}
	return
}
