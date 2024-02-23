package mdl

import (
	"strings"

	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/support/inflect"
	"github.com/ionous/errutil"
)

type fieldSet struct {
	kind   string
	fields []FieldInfo // a list of future fields
}

func (fs *fieldSet) writeFields(pen *Pen) (err error) {
	var cache fieldCache
	if kind, e := pen.findRequiredKind(fs.kind); e != nil {
		err = e
	} else if e := fs.rewriteImplicitAspects(pen, kind, &cache); e != nil {
		err = e
	} else if e := cache.writeFields(pen, kind, fs.fields, pen.addField); e != nil {
		err = e
	} else if e := fs.writeDefaultTraits(pen, kind, cache); e != nil {
		err = e
	}
	if err != nil {
		err = errutil.Fmt("%w in kind %q domain %q", err, fs.kind, pen.domain)
	}
	return
}

// rewrite object fields
func (fs *fieldSet) rewriteImplicitAspects(pen *Pen, kind kindInfo, cache *fieldCache) (err error) {
	if isObject := strings.HasSuffix(kind.fullpath(), pen.paths.kindsPath); isObject {
		for i := range fs.fields {
			field := &fs.fields[i]
			// rewrite Bool: "something" to an affinity with the opposite "not something" available.
			// i originally wanted to limit or force these into the format "is something"
			// but that screws with jess, sentences would have to be: "the noun is is something"
			if field.Affinity == affine.Bool && len(field.Class) == 0 {
				// default trait is the unset version
				defaultTrait := inflect.Join([]string{"not", field.Name})
				traits := []string{defaultTrait, field.Name}
				// rewrite bool fields as implicit aspects
				aspect := inflect.Join([]string{field.Name, "status"})
				cls, e := pen.addAspect(aspect, traits)
				if e := eatDuplicates(pen.warn, e); e != nil {
					err = e
					break
				} else {
					*field = FieldInfo{
						Name:     aspect,
						Class:    aspect,
						Affinity: affine.Text,
						Init:     field.Init,
					}
					cache.store(aspect, cls)
				}
			}
		}
	}
	return
}

// give aspect fields a provisional default
func (fs *fieldSet) writeDefaultTraits(pen *Pen, kind kindInfo, cache fieldCache) (err error) {
	if isObject := strings.HasSuffix(kind.fullpath(), pen.paths.kindsPath); isObject {
		for i := range fs.fields {
			field := &fs.fields[i]
			if field.isAspectLike() {
				aspect := cache[field.getClass()]
				if strings.HasSuffix(aspect.fullpath(), pen.paths.aspectPath) {
					if defaultTrait, e := pen.findDefaultTrait(aspect.class()); e != nil {
						err = e
						break
					} else if e := pen.addDefaultValue(kind, field.Name, ProvisionalAssignment{
						&assign.FromText{Value: &literal.TextValue{
							Value: defaultTrait,
						}}}); e != nil {
						err = e
						break
					}
				}
			}
		}
	}
	return
}
