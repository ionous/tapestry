package story

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/imp"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
	"github.com/ionous/errutil"
)

// Execute - called by the macro runtime during weave.
func (op *ActionDecl) Execute(macro rt.Runtime) error {
	return imp.StoryStatement(macro, op)
}

// Execute - actions generate pattern ephemera.
func (op *ActionDecl) PostImport(k *imp.Importer) (err error) {
	patterns := []string{op.Action.Str, op.Event.Str}
	targets := []string{"agent", "actor"}
	ancestors := []kindsOf.Kinds{kindsOf.Action, kindsOf.Event}
	for i, ancestor := range ancestors {
		tgt := targets[i]
		pattern := patterns[i]
		if e := k.AssertAncestor(pattern, ancestor.String()); e != nil {
			err = e
			break
		} else if e := k.AssertParam(pattern, tgt, tgt, affine.Text, nil); e != nil {
			err = e // ^ the first parameters is always (ex) "agent" of type "agent".
			break
		} else if extras, ok := op.ActionParams.Value.(FieldDefinition); !ok {
			err = errutil.New("unknown field type %T", op.ActionParams.Value)
			break
		} else if e := extras.DeclareField(func(name, class string, aff affine.Affinity, init assign.Assignment) error {
			return k.AssertParam(pattern, name, class, aff, init)
		}); e != nil {
			err = e
			break
		} else if i == 1 {
			if e := k.AssertResult(pattern, "success", "", affine.Bool, nil); e != nil {
				err = e
				break
			}
		}

	}
	return
}

const actionNoun = "noun"
const actionOtherNoun = "other_noun"

func (op *CommonAction) DeclareField(fn fieldType) error {
	return fn(actionNoun, op.Kind.Str, affine.Text, nil)
}

func (op *PairedAction) DeclareField(fn fieldType) (err error) {
	if e := fn(actionNoun, op.Kinds.Str, affine.Text, nil); e != nil {
		err = e
	} else if e := fn(actionOtherNoun, op.Kinds.Str, affine.Text, nil); e != nil {
		err = e
	}
	return
}

func (op *AbstractAction) DeclareField(fn fieldType) (none error) {
	return // no extra parameters
}
