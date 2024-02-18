package story

import (
	"errors"

	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/dl/jess"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/support/grok"
	"git.sr.ht/~ionous/tapestry/weave"
	"git.sr.ht/~ionous/tapestry/weave/mdl"
	"github.com/ionous/errutil"

	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
	"git.sr.ht/~ionous/tapestry/rt/safe"
)

// Execute - called by the macro runtime during weave.
func (op *DeclareStatement) Execute(macro rt.Runtime) error {
	return Weave(macro, op)
}

func (op *DeclareStatement) Weave(cat *weave.Catalog) error {
	return cat.Schedule(weave.RequireDependencies, func(w *weave.Weaver) (err error) {
		if text, e := safe.GetText(w, op.Text); e != nil {
			err = e
		} else if spans, e := grok.MakeSpans(text.String()); e != nil {
			err = e
		} else {
			// split each statement into its own evaluation
			// ( to break up interdependence )
			for _, temp := range spans {
				span := temp // pin otherwise the callback(s) all see the same last loop value
				if e := cat.Schedule(weave.RequireRules, func(w *weave.Weaver) (err error) {
					if i, e := w.MatchSpan(span); e != nil {
						err = errutil.Fmt("%w reading of %v b/c %v", mdl.Missing, span.String(), e)
					} else {
						err = i.Apply(jessAdapter{w, w.Pin()})
					}
					return
				}); e != nil {
					err = e
					break
				}
			}
		}
		return
	})
}

// fix... obviously.
type jessAdapter struct {
	w   *weave.Weaver
	pen *mdl.Pen
}

func (ja jessAdapter) AddKind(kind, ancestor string) error {
	return ja.pen.AddKind(kind, ancestor)
}
func (ja jessAdapter) AddKindTrait(kind, trait string) error {
	return ja.pen.AddKindTrait(kind, trait)
}
func (ja jessAdapter) AddNoun(short, long, kind string) error {
	return ja.pen.AddNoun(short, long, kind)
}
func (ja jessAdapter) AddNounAlias(noun, name string, rank int) error {
	return ja.pen.AddNounAlias(noun, name, rank)
}
func (ja jessAdapter) AddNounTrait(noun, trait string) error {
	return ja.w.AddInitialValue(ja.pen, noun, trait, truly())
}
func (ja jessAdapter) AddNounValue(noun, prop string, val rt.Assignment) error {
	return ja.w.AddInitialValue(ja.pen, noun, prop, val)
}
func (ja jessAdapter) GetClosestNoun(name string) (string, error) {
	return ja.w.GetClosestNoun(name)
}
func (ja jessAdapter) GetExactNoun(name string) (string, error) {
	return ja.pen.GetExactNoun(name)
}
func (ja jessAdapter) GetPlural(word string) string {
	return ja.w.PluralOf(word)
}
func (ja jessAdapter) GetSingular(word string) string {
	return ja.w.SingularOf(word)
}
func (ja jessAdapter) GetUniqueName(category string) string {
	return newCounter(ja.w.Catalog, category)
}

func (ja jessAdapter) Apply(macro jess.Macro, lhs, rhs []string) (err error) {
	if multiSrc, e := validSources(lhs, macro.Type); e != nil {
		err = e
	} else if multiTgt, e := validTargets(rhs, macro.Type); e != nil {
		err = e
	} else {
		src := namesToValues(lhs, multiSrc)
		tgt := namesToValues(rhs, multiTgt)
		if kind, e := ja.w.GetKindByName(macro.Name); e != nil {
			err = e
		} else if !kind.Implements(kindsOf.Macro.String()) {
			err = errutil.Fmt("expected %q to be a macro", kind.Name())
		} else if fieldCnt := kind.NumField(); fieldCnt < 2 {
			err = errutil.Fmt("expected macro %q to have at least two argument (not %d)", kind.Name(), fieldCnt)
		} else {
			args := []g.Value{src, tgt}
			if v, e := ja.w.Call(kind.Name(), affine.Text, nil, args); e != nil && !errors.Is(e, rt.NoResult) {
				err = e
			} else if v != nil {
				if msg := v.String(); len(msg) > 0 {
					err = errutil.Fmt("Declare statement: %s", msg)
				}
			}
		}
	}
	return
}

// validate that the number of parsed primary names is as expected
func validSources(ns []string, mtype grok.MacroType) (multi bool, err error) {
	switch mtype {
	case grok.Macro_PrimaryOnly, grok.Macro_ManyPrimary, grok.Macro_ManyMany:
		if cnt := len(ns); cnt == 0 {
			err = errutil.New("expected at least one source noun")
		}
		multi = true
	case grok.Macro_ManySecondary:
		if cnt := len(ns); cnt > 1 {
			err = errutil.New("expected exactly one noun")
		}
	default:
		err = errutil.New("invalid macro type")
	}
	return
}

// validate that the number of parsed secondary names is as expected
func validTargets(ns []string, mtype grok.MacroType) (multi bool, err error) {
	switch mtype {
	case grok.Macro_PrimaryOnly:
		if cnt := len(ns); cnt != 0 {
			err = errutil.New("didn't expect any target nouns")
		}
	case grok.Macro_ManyPrimary:
		if cnt := len(ns); cnt > 1 {
			err = errutil.New("expected at most one target noun")
		}
	case grok.Macro_ManySecondary, grok.Macro_ManyMany:
		// any number okay
		multi = true
	default:
		err = errutil.New("invalid macro type")
	}
	return
}

func namesToValues(names []string, multi bool) (ret g.Value) {
	if multi {
		ret = g.StringsOf(names)
	} else if len(names) > 0 {
		ret = g.StringOf(names[0])
	} else {
		ret = g.Empty
	}
	return
}
