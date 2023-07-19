package story

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/lang"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"git.sr.ht/~ionous/tapestry/support/grok"
	"git.sr.ht/~ionous/tapestry/weave"
	"git.sr.ht/~ionous/tapestry/weave/mdl"
	"github.com/ionous/errutil"
	"github.com/ionous/sliceOf"
)

// Execute - called by the macro runtime during weave.
func (op *DefineKinds) Execute(macro rt.Runtime) error {
	return Weave(macro, op)
}

// ex. "cats are a kind of animal"
func (op *DefineKinds) Weave(cat *weave.Catalog) error {
	return cat.Schedule(weave.RequirePlurals, func(w *weave.Weaver) (err error) {
		if kinds, e := safe.GetTextList(w, op.Kinds); e != nil {
			err = e
		} else if ancestor, e := safe.GetText(w, op.Ancestor); e != nil {
			err = e
		} else {
			pen := w.Pin()
			ancestor := lang.Normalize(ancestor.String())
			for _, kind := range kinds.Strings() {
				// tbd: are the determiners of kinds useful for anything?
				if kind, e := grok.StripArticle(kind); e != nil {
					err = errutil.Append(err, e)
				} else {
					kind := lang.Normalize(kind)
					if e := pen.AddKind(kind, ancestor); e != nil {
						err = errutil.Append(err, e)
					}
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
func (op *DefineFields) Weave(cat *weave.Catalog) (err error) {
	return cat.Schedule(weave.RequirePlurals, func(w *weave.Weaver) (err error) {
		if kind, e := safe.GetText(cat.Runtime(), op.Kind); e != nil {
			err = e
		} else {
			pen := w.Pin()
			fields := mdl.NewFieldBuilder(kind.String())
			for _, field := range op.Fields {
				// bools here become implicit aspects.
				// ( vs. bool pattern vars which stay bools -- see reduceProps )
				if el := field.FieldInfo(); el.Affinity != affine.Bool {
					fields.AddField(el)
				} else {
					// tbd: would it be nicer to support single trait kinds?
					// not_aspect would instead be: Not{IsTrait{PositiveName}}
					if e := pen.AddAspect(el.Name, sliceOf.String("not "+el.Name, "is "+el.Name)); e != nil {
						err = errutil.Append(err, e)
					} else {
						fields.AddAspect(el.Name)
					}
				}
			}
			if err == nil {
				err = pen.AddFields(fields.Fields)
			}
		}
		return
	})
}
