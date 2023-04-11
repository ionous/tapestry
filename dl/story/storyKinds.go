package story

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/imp"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"github.com/ionous/errutil"
)

// Execute - called by the macro runtime during weave.
func (op *DefineKinds) Execute(macro rt.Runtime) error {
	return imp.StoryStatement(macro, op)
}

// ex. "cats are a kind of animal"
func (op *DefineKinds) PostImport(k *imp.Importer) (err error) {
	// FIX: macro runtime
	if kinds, e := safe.GetTextList(k, op.Kinds); e != nil {
		err = e
	} else if ancestor, e := safe.GetText(k, op.Ancestor); e != nil {
		err = e
	} else {
		for _, kind := range kinds.Strings() {
			if e := k.AssertAncestor(kind, ancestor.String()); e != nil {
				err = errutil.Append(err, e)
			}
		}
	}
	return
}

// Execute - called by the macro runtime during weave.
func (op *DefineFields) Execute(macro rt.Runtime) error {
	return imp.StoryStatement(macro, op)
}

// ex. cats have some text called breed.
// ex. horses have an aspect called speed.
func (op *DefineFields) PostImport(k *imp.Importer) (err error) {
	if len(op.Fields) == 0 {
		// log or something?
	} else if kind, e := safe.GetText(k, op.Kind); e != nil {
		err = e
	} else if len(op.Fields) > 0 {
		kind := kind.String()
		for _, el := range op.Fields {
			// handle every other than bools....
			if aspect, ok := el.(*BoolField); !ok {
				err = el.DeclareField(func(name, class string, aff affine.Affinity, init assign.Assignment) (err error) {
					return k.AssertField(kind, name, class, aff, init)
				})
			} else {
				// bools here become implicit aspects.
				// ( vs. bool pattern vars which stay bools -- see reduceProps )
				// first: add the aspect
				aspect := aspect.Name
				traits := []string{"not_" + aspect, "is_" + aspect}

				// fix: future: it'd be nicer to support single trait kinds
				// not_aspect would instead be: Not{IsTrait{PositiveName}}
				if k.AssertAspectTraits(aspect, traits); e != nil {
					err = errutil.Append(err, e)
				} else if e := k.AssertField(kind, aspect, aspect, affine.Text, nil); e != nil {
					err = errutil.Append(err, e)
				}
			}
		}
	}
	return
}
