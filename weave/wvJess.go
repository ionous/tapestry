package weave

import (
	"errors"
	"strconv"

	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/dl/jess"
	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/rt"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
	"git.sr.ht/~ionous/tapestry/support/jessdb"
	"git.sr.ht/~ionous/tapestry/weave/mdl"
	"github.com/ionous/errutil"
)

// fix... obviously.
// mostly its the same as Pen, but there are some overrides, and rename *sigh*
type jessAdapter struct {
	w *Weaver
	*mdl.Pen
}

// PostProcess schedules a jess command for later handling
func (ja jessAdapter) PostProcess(pri jess.Priority, post jess.Process) error {
	var phase Phase
	switch pri {
	case jess.NounSettings:
		phase = RequireNouns
	case jess.Understandings:
		phase = RequireAll
	default:
		panic("unexpected priority")
	}
	return ja.w.Catalog.Schedule(phase, func(w *Weaver) error {
		q := jessdb.MakeQuery(w.Catalog.Modeler, w.Domain)
		return post(q, ja)
	})
}

func (ja jessAdapter) AddNounTrait(noun, trait string) error {
	return ja.w.AddNounValue(ja.Pen, noun, trait, truly())
}
func (ja jessAdapter) AddNounValue(noun, prop string, val rt.Assignment) error {
	return ja.w.AddNounValue(ja.Pen, noun, prop, val)
}
func (ja jessAdapter) AddNounPath(noun string, path []string, val literal.LiteralValue) error {
	return ja.w.AddNounPath(ja.Pen, noun, path, val)
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
func (ja jessAdapter) GetUniqueName(category string) string {
	return newCounter(ja.w.Catalog, category)
}

func (ja jessAdapter) Apply(macro mdl.Macro, lhs, rhs []string) (err error) {
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
func validSources(ns []string, mtype mdl.MacroType) (multi bool, err error) {
	switch mtype {
	case mdl.Macro_PrimaryOnly, mdl.Macro_ManyPrimary, mdl.Macro_ManyMany:
		if cnt := len(ns); cnt == 0 {
			err = errutil.New("expected at least one source noun")
		}
		multi = true
	case mdl.Macro_ManySecondary:
		if cnt := len(ns); cnt > 1 {
			err = errutil.New("expected exactly one noun")
		}
	default:
		err = errutil.New("invalid macro type")
	}
	return
}

// validate that the number of parsed secondary names is as expected
func validTargets(ns []string, mtype mdl.MacroType) (multi bool, err error) {
	switch mtype {
	case mdl.Macro_PrimaryOnly:
		if cnt := len(ns); cnt != 0 {
			err = errutil.New("didn't expect any target nouns")
		}
	case mdl.Macro_ManyPrimary:
		if cnt := len(ns); cnt > 1 {
			err = errutil.New("expected at most one target noun")
		}
	case mdl.Macro_ManySecondary, mdl.Macro_ManyMany:
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

func truly() rt.Assignment {
	return &assign.FromBool{
		Value: &literal.BoolValue{Value: true},
	}
}

func newCounter(cat *Catalog, name string) (ret string) {
	next := cat.Env.Inc(name, 1)
	return name + "-" + strconv.Itoa(next)
}
