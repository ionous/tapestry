package story

import (
	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
	"git.sr.ht/~ionous/tapestry/rt/meta"
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"git.sr.ht/~ionous/tapestry/support/inflect"
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

func (op *DefineNounKind) Weave(cat *weave.Catalog) error {
	return cat.Schedule(weaver.NounPhase, func(w weaver.Weaves, run rt.Runtime) (err error) {
		if noun, e := safe.GetText(run, op.NounName); e != nil {
			err = e
		} else if kind, e := safe.GetText(run, op.KindName); e != nil {
			err = e
		} else {
			n := inflect.Normalize(noun.String())
			k := inflect.Normalize(kind.String())
			_, err = mdl.AddNamedNoun(w, n, k)
		}
		return
	})
}

func (op *DefineNounStates) Weave(cat *weave.Catalog) error {
	return cat.Schedule(weaver.ValuePhase, func(w weaver.Weaves, run rt.Runtime) (err error) {
		if name, e := safe.GetText(run, op.NounName); e != nil {
			err = e
		} else if traits, e := safe.GetTextList(run, op.StateNames); e != nil {
			err = e
		} else if noun, e := run.GetField(meta.ObjectId, name.String()); e != nil {
			err = e
		} else {
			for _, t := range traits.Strings() {
				t := inflect.Normalize(t)
				if e := w.AddNounValue(noun.String(), t, truly()); e != nil {
					err = e
					break
				}
			}
		}
		return
	})
}

// ex. The description of the nets is xxx
func (op *DefineNounValue) Weave(cat *weave.Catalog) error {
	return cat.Schedule(weaver.ValuePhase, func(w weaver.Weaves, run rt.Runtime) (err error) {
		if name, e := safe.GetText(run, op.NounName); e != nil {
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
				field := inflect.Normalize(field.String())
				if noun, e := run.GetField(meta.ObjectId, name.String()); e != nil {
					err = e
				} else if e := w.AddNounValue(noun.String(), field, value); e != nil {
					err = e
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
