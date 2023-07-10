package story

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"git.sr.ht/~ionous/tapestry/weave"
	"git.sr.ht/~ionous/tapestry/weave/assert"
	"github.com/ionous/errutil"
)

// Execute - called by the macro runtime during weave.
func (op *DefineKinds) Execute(macro rt.Runtime) error {
	return Weave(macro, op)
}

// ex. "cats are a kind of animal"
func (op *DefineKinds) Weave(cat *weave.Catalog) error {
	return cat.Schedule(assert.RequirePlurals, func(w *weave.Weaver) (err error) {
		if kinds, e := safe.GetTextList(w, op.Kinds); e != nil {
			err = e
		} else if ancestor, e := safe.GetText(w, op.Ancestor); e != nil {
			err = e
		} else {
			for _, kind := range kinds.Strings() {
				if e := cat.AssertAncestor(kind, ancestor.String()); e != nil {
					err = errutil.Append(err, e)
				}
			}
		}
		return
	})
}

// Execute - called by the macro runtime during weave.
func (op *DefineFields) Execute(macro rt.Runtime) error {
	return Weave(macro, op)
}

// ex. cats have some text called breed.
// ex. horses have an aspect called speed.
func (op *DefineFields) Weave(cat *weave.Catalog) error {
	return cat.Schedule(assert.RequireDeterminers, func(w *weave.Weaver) (err error) {
		if len(op.Fields) == 0 {
			// log or something?
		} else if kind, e := safe.GetText(w, op.Kind); e != nil {
			err = e
		} else if len(op.Fields) > 0 {
			kind := kind.String()
			for _, el := range op.Fields {
				// handle every other than bools....
				if aspect, ok := el.(*BoolField); !ok {
					err = el.DeclareField(func(name, class string, aff affine.Affinity, init assign.Assignment) (err error) {
						return cat.AssertField(kind, name, class, aff, init)
					})
				} else {
					// bools here become implicit aspects.
					// ( vs. bool pattern vars which stay bools -- see reduceProps )
					// first: add the aspect
					aspect := aspect.Name
					traits := []string{"not " + aspect, "is " + aspect}

					// fix: future: it'd be nicer to support single trait kinds
					// not_aspect would instead be: Not{IsTrait{PositiveName}}
					if cat.AssertAspectTraits(aspect, traits); e != nil {
						err = errutil.Append(err, e)
					} else if e := cat.AssertField(kind, aspect, aspect, affine.Text, nil); e != nil {
						err = errutil.Append(err, e)
					}
				}
			}
		}
		return
	})
}
