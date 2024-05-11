package story

import (
	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
	"git.sr.ht/~ionous/tapestry/rt/meta"
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"git.sr.ht/~ionous/tapestry/support/inflect"
	"git.sr.ht/~ionous/tapestry/support/match"
	"git.sr.ht/~ionous/tapestry/weave"
	"git.sr.ht/~ionous/tapestry/weave/mdl"
	"git.sr.ht/~ionous/tapestry/weave/weaver"
	"github.com/ionous/errutil"
)

// (the) colors are red, blue, or green.
func (op *DefineState) Weave(cat *weave.Catalog) error {
	return cat.Schedule(weaver.AncestryPhase, func(w weaver.Weaves, run rt.Runtime) (err error) {
		if aspect, e := safe.GetText(run, op.Aspect); e != nil {
			err = e
		} else if traits, e := safe.GetTextList(run, op.Traits); e != nil {
			err = e
		} else {
			aspect, traits := inflect.Normalize(aspect.String()), traits.Strings()
			for i, t := range traits {
				traits[i] = inflect.Normalize(t)
			}
			if e := w.AddKind(aspect, kindsOf.Aspect.String()); e != nil {
				err = e
			} else {
				// we could wait till AncestryPhase, but now is fine too.
				err = w.AddAspectTraits(aspect, traits)
			}
		}
		return
	})
}

func (op *DefineNounTraits) Weave(cat *weave.Catalog) error {
	return cat.Schedule(weaver.ValuePhase, func(w weaver.Weaves, run rt.Runtime) (err error) {
		if nouns, e := safe.GetTextList(run, op.Nouns); e != nil {
			err = e
		} else if traits, e := safe.GetTextList(run, op.Traits); e != nil {
			err = e
		} else if traits := traits.Strings(); len(traits) > 0 {
			names := nouns.Strings()
			for _, t := range traits {
				t := inflect.Normalize(t)
				for _, n := range names {
					if e := w.AddNounValue(n, t, truly()); e != nil {
						err = errutil.Append(err, e)
						break // out of the traits to the next noun
					}
				}
			}
		}
		return
	})
}

func (op *DefineNouns) Weave(cat *weave.Catalog) error {
	return cat.Schedule(weaver.NounPhase, func(w weaver.Weaves, run rt.Runtime) (err error) {
		if nouns, e := safe.GetTextList(run, op.Nouns); e != nil {
			err = e
		} else if kind, e := safe.GetText(run, op.Kind); e != nil {
			err = e
		} else {
			names := nouns.Strings()
			if kind := kind.String(); len(kind) > 0 {
				kind := match.StripArticle(kind)
				for _, noun := range names {
					noun := match.StripArticle(noun)
					if _, e := mdl.AddNamedNoun(w, noun, kind); e != nil {
						err = errutil.Append(err, e)
					}
				}
			}
		}
		return
	})
}

// ex. The description of the nets is xxx
func (op *DefineValue) Weave(cat *weave.Catalog) error {
	return cat.Schedule(weaver.ValuePhase, func(w weaver.Weaves, run rt.Runtime) (err error) {
		if nouns, e := safe.GetTextList(run, op.Nouns); e != nil {
			err = e
		} else if field, e := safe.GetText(run, op.FieldName); e != nil {
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
				subjects := nouns.Strings()
				field := field.String()
				for _, noun := range subjects {
					name := match.StripArticle(noun)
					if noun, e := run.GetField(meta.ObjectId, name); e != nil {
						err = errutil.Append(err, e)
					} else if e := w.AddNounValue(noun.String(), field, value); e != nil {
						err = errutil.Append(err, e)
					}
				}
			}
		}
		return
	})
}

func (op *DefineRelatives) Weave(cat *weave.Catalog) error {
	return cat.Schedule(weaver.ConnectionPhase, func(w weaver.Weaves, run rt.Runtime) (err error) {
		if rel, e := safe.GetText(run, op.Relation); e != nil {
			err = e
		} else if nouns, e := safe.GetTextList(run, op.Nouns); e != nil {
			err = e
		} else if otherNouns, e := safe.GetTextList(run, op.OtherNouns); e != nil {
			err = e
		} else {
			err = defineRelatives(w, run, rel.String(), nouns.Strings(), otherNouns.Strings())
		}
		return
	})
}

func (op *DefineOtherRelatives) Weave(cat *weave.Catalog) error {
	return cat.Schedule(weaver.ConnectionPhase, func(w weaver.Weaves, run rt.Runtime) (err error) {
		if rel, e := safe.GetText(run, op.Relation); e != nil {
			err = e
		} else if nouns, e := safe.GetTextList(run, op.Nouns); e != nil {
			err = e
		} else if otherNouns, e := safe.GetTextList(run, op.OtherNouns); e != nil {
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
			err = defineRelatives(w, run, rel.String(), otherNouns.Strings(), nouns.Strings())
		}
		return
	})
}

func defineRelatives(w weaver.Weaves, run rt.Runtime, relation string, nouns, otherNouns []string) (err error) {
	rel := inflect.Normalize(relation)
	for _, one := range nouns {
		if a, e := run.GetField(meta.ObjectId, one); e != nil {
			err = errutil.Append(err, e)
		} else {
			a := a.String()
			for _, other := range otherNouns {
				if b, e := run.GetField(meta.ObjectId, other); e != nil {
					err = errutil.Append(err, e)
				} else {
					b := b.String()
					if e := w.AddNounPair(rel, a, b); e != nil {
						err = errutil.Append(err, e)
					}
				}
			}
		}
	}
	return
}
