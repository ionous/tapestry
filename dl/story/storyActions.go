package story

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"git.sr.ht/~ionous/tapestry/weave"
	"git.sr.ht/~ionous/tapestry/weave/assert"
	"github.com/ionous/errutil"
)

// Execute - called by the macro runtime during weave.
func (op *ActionDecl) Execute(macro rt.Runtime) error {
	return Weave(macro, op)
}

// Execute - actions generate pattern ephemera.
func (op *ActionDecl) Weave(cat *weave.Catalog) (err error) {
	return cat.Schedule(assert.RequireAncestry, func(w *weave.Weaver) (err error) {
		patterns := []string{op.Action.Str, op.Event.Str}
		targets := []string{"agent", "actor"}
		ancestors := []kindsOf.Kinds{kindsOf.Action, kindsOf.Event}
		for i, ancestor := range ancestors {
			tgt := targets[i]
			pattern := patterns[i]
			if e := cat.AssertAncestor(pattern, ancestor.String()); e != nil {
				err = e
				break
			} else if e := cat.AssertParam(pattern, tgt, tgt, affine.Text, nil); e != nil {
				err = e // ^ the first parameters is always (ex) "agent" of type "agent".
				break
			} else if e := op.defineExtras(w, cat, pattern); e != nil {
				err = e
				break
			} else if i == 1 {
				if e := cat.AssertResult(pattern, "success", "", affine.Bool, nil); e != nil {
					err = e
					break
				}
			}
		}
		return
	})
}

const actionNoun = "noun"
const actionOtherNoun = "other noun"

func (op *ActionDecl) defineExtras(run rt.Runtime, k assert.Assertions, pattern string) (err error) {
	extras := op.ActionParams.Value
	switch p := extras.(type) {
	case *CommonAction:
		if kind, e := safe.GetText(run, p.Kind); e != nil {
			err = e
		} else {
			err = k.AssertParam(pattern, actionNoun, kind.String(), affine.Text, nil)
		}
	case *PairedAction:
		if kinds, e := safe.GetText(run, p.Kinds); e != nil {
			err = e
		} else if e := k.AssertParam(pattern, actionNoun, kinds.String(), affine.Text, nil); e != nil {
			err = e
		} else {
			err = k.AssertParam(pattern, actionOtherNoun, kinds.String(), affine.Text, nil)
		}
	case *AbstractAction:
		// no extra parameters
	default:
		err = errutil.Fmt("unknown action %T", extras)
	}
	return
}
