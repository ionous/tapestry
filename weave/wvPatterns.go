package weave

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
	"git.sr.ht/~ionous/tapestry/weave/assert"
	"github.com/ionous/errutil"
)

// writes a definition of kindName?args=arg1,arg2,arg3
func (cat *Catalog) AssertParam(kind, field, class string, aff affine.Affinity, init assign.Assignment) error {
	return cat.Schedule(assert.RequireAncestry, func(ctx *Weaver) error {
		d, at := ctx.d, ctx.at
		return addField(ctx, kind, field, class, func(kind, field, class string) (err error) {
			if init != nil {
				err = errutil.New("parameters don't currently support initial values")
			} else if e := ctx.d.addKind(kind, kindsOf.Pattern.String(), at); e != nil {
				err = e
			} else {
				err = cat.writer.Parameter(d.name, kind, field, aff, class, at)
			}
			return
		})
	})
}

func (cat *Catalog) AssertResult(kind, field, class string, aff affine.Affinity, init assign.Assignment) error {
	return cat.Schedule(assert.RequireParameters, func(ctx *Weaver) error {
		d, at := ctx.d, ctx.at
		return addField(ctx, kind, field, class, func(kind, field, class string) (err error) {
			if init != nil {
				err = errutil.New("return values don't currently support initial values")
			} else if e := ctx.d.addKind(kind, kindsOf.Pattern.String(), at); e != nil {
				err = e
			} else {
				err = cat.writer.Result(d.name, kind, field, aff, class, at)
			}
			return
		})
	})
}

func (cat *Catalog) AssertField(kind, field, class string, aff affine.Affinity, init assign.Assignment) error {
	return cat.Schedule(assert.RequireResults, func(ctx *Weaver) error {
		d, at := ctx.d, ctx.at
		return addField(ctx, kind, field, class, func(kind, field, class string) (err error) {
			// shortcut: if we specify a field name for a record and no class, we'll expect the class to be the name.
			if len(class) == 0 && isRecordAffinity(aff) {
				class = field
			}
			if e := cat.writer.Member(d.name, kind, field, aff, class, at); e != nil {
				err = e
			} else if init != nil {
				return cat.writer.Default(d.name, kind, field, init)
			}
			return
		})
	})
}

// ugly way to normalize field names
func addField(ctx *Weaver, kindName, fieldName, fieldClass string,
	addField func(k, f, c string) error) (err error) {
	d := ctx.d
	if _, kind := d.UniformDeterminer(kindName); len(kind) == 0 {
		err = InvalidString(kindName)
	} else if field, ok := UniformString(fieldName); !ok {
		err = InvalidString(fieldName)
	} else if class, ok := UniformString(fieldClass); !ok && len(fieldClass) > 0 {
		err = InvalidString(fieldClass)
	} else {
		kind, class := d.singularize(kind), d.singularize(class)
		err = addField(kind, field, class)
	}
	return
}
