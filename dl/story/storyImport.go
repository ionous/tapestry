package story

import (
	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"git.sr.ht/~ionous/tapestry/support/inflect"
	"git.sr.ht/~ionous/tapestry/support/match"
	"git.sr.ht/~ionous/tapestry/weave"
	"git.sr.ht/~ionous/tapestry/weave/mdl"
	"github.com/ionous/errutil"
)

// Execute - called by the macro runtime during weave.
func (op *DefineAspect) Execute(macro rt.Runtime) error {
	return Weave(macro, op)
}

// (the) colors are red, blue, or green.
func (op *DefineAspect) Weave(cat *weave.Catalog) error {
	return cat.Schedule(weave.LanguagePhase, func(w *weave.Weaver) (err error) {
		if aspect, e := safe.GetText(w, op.Aspect); e != nil {
			err = e
		} else if traits, e := safe.GetTextList(w, op.Traits); e != nil {
			err = e
		} else {
			aspect, traits := inflect.Normalize(aspect.String()), traits.Strings()
			for i, t := range traits {
				traits[i] = inflect.Normalize(t)
			}
			pen := w.Pin()
			if e := pen.AddKind(aspect, kindsOf.Aspect.String()); e != nil {
				err = e
			} else {
				err = pen.AddAspectTraits(aspect, traits)
			}
		}
		return
	})
}

// Execute - called by the macro runtime during weave.
func (op *DefineNounTraits) Execute(macro rt.Runtime) error {
	return Weave(macro, op)
}

func (op *DefineNounTraits) Weave(cat *weave.Catalog) error {
	return cat.Schedule(weave.LanguagePhase, func(w *weave.Weaver) (err error) {
		if nouns, e := safe.GetTextList(w, op.Nouns); e != nil {
			err = e
		} else if traits, e := safe.GetTextList(w, op.Traits); e != nil {
			err = e
		} else if traits := traits.Strings(); len(traits) > 0 {
			pen := w.Pin()
			names := nouns.Strings()
			for _, t := range traits {
				t := inflect.Normalize(t)
				for _, n := range names {
					if e := w.AddNounValue(pen, n, t, truly()); e != nil {
						err = errutil.Append(err, e)
						break // out of the traits to the next noun
					}
				}
			}
		}
		return
	})
}

// Execute - called by the macro runtime during weave.
func (op *DefinePhrase) Execute(macro rt.Runtime) error {
	return Weave(macro, op)
}

func (op *DefinePhrase) Weave(cat *weave.Catalog) error {
	return cat.Schedule(weave.AncestryPhase, func(w *weave.Weaver) (err error) {
		if phrase, e := safe.GetText(w, op.Phrase); e != nil {
			err = e
		} else if macro, e := safe.GetText(w, op.Macro); e != nil {
			err = e
		} else if rev, e := safe.GetOptionalBool(w, op.Reversed, false); e != nil {
			err = e
		} else if macro := inflect.Normalize(macro.String()); len(macro) == 0 {
			err = errutil.New("missing macro name")
		} else {
			err = w.Pin().AddPhrase(macro, phrase.String(), rev.Bool())
		}
		return
	})
}

// Execute - called by the macro runtime during weave.
func (op *DefineNouns) Execute(macro rt.Runtime) error {
	return Weave(macro, op)
}

func (op *DefineNouns) Weave(cat *weave.Catalog) error {
	return cat.Schedule(weave.LanguagePhase, func(w *weave.Weaver) (err error) {
		if nouns, e := safe.GetTextList(w, op.Nouns); e != nil {
			err = e
		} else if kind, e := safe.GetText(w, op.Kind); e != nil {
			err = e
		} else {
			names := nouns.Strings()
			if kind := kind.String(); len(kind) > 0 {
				pen := w.Pin()
				kind := match.StripArticle(kind)
				for _, noun := range names {
					noun := match.StripArticle(noun)
					if _, e := mdl.AddNamedNoun(pen, noun, kind); e != nil {
						err = errutil.Append(err, e)
					}
				}
			}
		}
		return
	})
}

// Execute - called by the macro runtime during weave.
func (op *DefineValue) Execute(macro rt.Runtime) error {
	return Weave(macro, op)
}

// ex. The description of the nets is xxx
func (op *DefineValue) Weave(cat *weave.Catalog) error {
	return cat.Schedule(weave.ValuePhase, func(w *weave.Weaver) (err error) {
		if nouns, e := safe.GetTextList(w, op.Nouns); e != nil {
			err = e
		} else if field, e := safe.GetText(w, op.FieldName); e != nil {
			err = e
		} else {
			// try to convert from literal templates ( if any )
			value := op.Value
			switch wrapper := op.Value.(type) {
			case *assign.FromText:
				switch text := wrapper.Value.(type) {
				case *literal.TextValue:
					value, err = convertTextAssignment(text.Value)
				}
			}
			if err == nil {
				pen := w.Pin()
				subjects := nouns.Strings()
				field := field.String()
				for _, noun := range subjects {
					name := match.StripArticle(noun)
					if noun, e := pen.GetClosestNoun(inflect.Normalize(name)); e != nil {
						err = errutil.Append(err, e)
					} else if e := w.AddNounValue(pen, noun, field, value); e != nil {
						err = errutil.Append(err, e)
					}
				}
			}
		}
		return
	})
}

// Execute - called by the macro runtime during weave.
func (op *DefineRelatives) Execute(macro rt.Runtime) error {
	return Weave(macro, op)
}

func (op *DefineRelatives) Weave(cat *weave.Catalog) error {
	return cat.Schedule(weave.LanguagePhase, func(w *weave.Weaver) (err error) {
		if rel, e := safe.GetText(w, op.Relation); e != nil {
			err = e
		} else if nouns, e := safe.GetTextList(w, op.Nouns); e != nil {
			err = e
		} else if otherNouns, e := safe.GetTextList(w, op.OtherNouns); e != nil {
			err = e
		} else {
			err = defineRelatives(w, rel.String(), nouns.Strings(), otherNouns.Strings())
		}
		return
	})
}

// Execute - called by the macro runtime during weave.
func (op *DefineOtherRelatives) Execute(macro rt.Runtime) error {
	return Weave(macro, op)
}

func (op *DefineOtherRelatives) Weave(cat *weave.Catalog) error {
	return cat.Schedule(weave.LanguagePhase, func(w *weave.Weaver) (err error) {
		if rel, e := safe.GetText(w, op.Relation); e != nil {
			err = e
		} else if nouns, e := safe.GetTextList(w, op.Nouns); e != nil {
			err = e
		} else if otherNouns, e := safe.GetTextList(w, op.OtherNouns); e != nil {
			err = e
		} else {
			// fix: for nearly every statement; after we resolve the names -- we'd want to re-schedule
			// because if there's anything mssing; we will lose the context ( ex. pattern call stack )
			// needed to re-resolve the names. capturing the context might work, but re-issuing the schedule might make sense too
			// ( ex. could make it RequireNames instead of requirePlurals, etc. )
			// even more interesting are "partial resolutions" --
			// this is where the promise api (weave.res) would come in handy --
			// resolve the rel, resolve the nouns promise all
			// then define the relation.
			err = defineRelatives(w, rel.String(), otherNouns.Strings(), nouns.Strings())
		}
		return
	})
}

func defineRelatives(w *weave.Weaver, rel string, nouns, otherNouns []string) (err error) {
	pen, rel := w.Pin(), inflect.Normalize(rel)
	for _, one := range nouns {
		if a, e := w.GetClosestNoun(inflect.Normalize(one)); e != nil {
			err = errutil.Append(err, e)
		} else {
			for _, other := range otherNouns {
				if b, e := w.GetClosestNoun(inflect.Normalize(other)); e != nil {
					err = errutil.Append(err, e)
				} else {
					if e := pen.AddNounPair(rel, a, b); e != nil {
						err = errutil.Append(err, e)
					}
				}
			}
		}
	}
	return
}
