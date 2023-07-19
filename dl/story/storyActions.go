package story

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/lang"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"git.sr.ht/~ionous/tapestry/weave"
	"git.sr.ht/~ionous/tapestry/weave/mdl"
	"github.com/ionous/errutil"
)

// Execute - called by the macro runtime during weave.
func (op *ActionDecl) Execute(macro rt.Runtime) error {
	return Weave(macro, op)
}

// Execute - actions generate pattern ephemera.
func (op *ActionDecl) Weave(cat *weave.Catalog) (err error) {
	return cat.Schedule(weave.RequireAncestry, func(w *weave.Weaver) (err error) {
		patterns := []string{op.Action.Str, op.Event.Str}
		targets := []string{"agent", "actor"}
		ancestors := []kindsOf.Kinds{kindsOf.Action, kindsOf.Event}
		for i, target := range targets {
			pb := mdl.NewPatternSubtype(patterns[i], ancestors[i])
			// the first parameter is always (ex) "agent" of type "agent":
			pb.AddField(mdl.PatternParameters, mdl.FieldInfo{
				Name:     target,
				Class:    target,
				Affinity: affine.Text,
			})
			//
			if e := op.defineExtras(cat.Runtime(), pb); e != nil {
				err = e
				break
			} else if i == 1 {
				pb.AddField(mdl.PatternResults, mdl.FieldInfo{
					Name:     "success",
					Affinity: affine.Bool,
				})
			}
			if e := w.Pin().AddPattern(pb.Pattern); e != nil {
				err = e
				break
			}
		}
		return
	})
}

const actionNoun = "noun"
const actionOtherNoun = "other noun"

func (op *ActionDecl) defineExtras(run rt.Runtime, pb *mdl.PatternBuilder) (err error) {
	extras := op.ActionParams.Value
	switch p := extras.(type) {
	case *CommonAction:
		if kind, e := safe.GetText(run, p.Kind); e != nil {
			err = e
		} else {
			cls := lang.Normalize(kind.String())
			pb.AddField(mdl.PatternParameters, mdl.FieldInfo{
				Name:     actionNoun,
				Class:    cls,
				Affinity: affine.Text,
			})
		}
	case *PairedAction:
		if kinds, e := safe.GetText(run, p.Kinds); e != nil {
			err = e
		} else {
			cls := lang.Normalize(kinds.String())
			pb.AddField(mdl.PatternParameters, mdl.FieldInfo{
				Name:     actionNoun,
				Class:    cls,
				Affinity: affine.Text,
			})
			pb.AddField(mdl.PatternParameters, mdl.FieldInfo{
				Name:     actionOtherNoun,
				Class:    cls,
				Affinity: affine.Text,
			})
		}
	case *AbstractAction:
		// no extra parameters
	default:
		err = errutil.Fmt("unknown action %T", extras)
	}
	return
}
