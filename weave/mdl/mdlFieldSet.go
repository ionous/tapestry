package mdl

import (
	"strings"

	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/dl/call"
	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
	"git.sr.ht/~ionous/tapestry/support/inflect"
	"fmt"
)

type fieldSet struct {
	kind   kindInfo
	fields []FieldInfo // a list of future fields
	cache  fieldCache
}

func (pen *Pen) writeFields(kind string, fields []FieldInfo) (err error) {
	if kind, e := pen.findRequiredKind(kind); e != nil {
		err = e
	} else {
		var cache fieldCache
		if e := cache.precache(pen, fields); e != nil {
			err = e
		} else {
			fs := fieldSet{kind, fields, cache}
			if e := fs.rewriteBooleans(pen); e != nil {
				err = e
			} else if e := fs.addFields(pen, pen.addField); e != nil {
				err = e
			} else if e := fs.writeDefaultTraits(pen); e != nil {
				err = e
			}
		}
	}
	if err != nil {
		err = fmt.Errorf("%w in kind %q domain %q", err, kind, pen.domain)
	}
	return
}

// rewrite object bool fields into two traits and an implicit aspect.
// ex. "bool portable" -> "portable status" aspect with "portable" and "not portable" states.
// the default value is the "not" version.
func (fs *fieldSet) rewriteBooleans(pen *Pen) (err error) {
	if isObject := strings.HasSuffix(fs.kind.fullpath(), pen.getPath(kindsOf.Kind)); isObject {
		for i, field := range fs.fields {
			// rewrite Bool: "something" to an affinity with the opposite "not something"
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
					fs.fields[i] = FieldInfo{
						Name:     aspect,
						Class:    aspect,
						Affinity: affine.Text,
						Init:     field.Init,
					}
					fs.cache.store(aspect, cls)
				}
			}
		}
	}
	return
}

type fieldHandler func(kid, cls kindInfo, field string, aff affine.Affinity) error

func (fs *fieldSet) addFields(pen *Pen, call fieldHandler) (err error) {
	for _, f := range fs.fields {
		if e := f.validate(); e != nil {
			err = e
		} else if cls, e := fs.cache.getClass(pen, f); e != nil {
			err = e
		} else {
			e := call(fs.kind, cls, f.Name, f.Affinity)
			if e := eatDuplicates(pen.warn, e); e != nil {
				err = e
			} else if f.Init != nil {
				e := pen.addDefaultValue(fs.kind, f.Name, f.Init)
				if e := eatDuplicates(pen.warn, e); e != nil {
					err = e
				}
			}
		}
		if err != nil {
			err = fmt.Errorf("%w trying to write field %q", err, f.Name)
			break
		}

	}
	return
}

// when adding fields of type aspect to a kind,
// set a provisional default trait
func (fs *fieldSet) writeDefaultTraits(pen *Pen) (err error) {
	kind := fs.kind
	if isObject := strings.HasSuffix(kind.fullpath(), pen.getPath(kindsOf.Kind)); isObject {
		for _, field := range fs.fields {
			if field.isAspectLike() {
				aspect := fs.cache[field.getDefaultClass()]
				if strings.HasSuffix(aspect.fullpath(), pen.getPath(kindsOf.Aspect)) {
					if defaultTrait, e := pen.findDefaultTrait(aspect.class()); e != nil {
						err = e
						break
					} else if e := pen.addDefaultValue(kind, field.Name, ProvisionalAssignment{
						&call.FromText{Value: &literal.TextValue{
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
