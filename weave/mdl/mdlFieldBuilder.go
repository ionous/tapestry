package mdl

import (
	"strings"

	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/lang"
	"git.sr.ht/~ionous/tapestry/rt"
	"github.com/ionous/errutil"
)

type FieldInfo struct {
	Name, Class string
	Affinity    affine.Affinity
	Init        rt.Assignment
}

type FieldBuilder struct {
	Fields
}

type Fields struct {
	kind string
	fieldSet
}

// supports enough slot for patterns
// see PatternBuilder.
type fieldSet struct {
	fields [NumFieldTypes][]FieldInfo
}

func NewFieldBuilder(kind string) *FieldBuilder {
	return &FieldBuilder{Fields: Fields{
		// tbd: feels like it'd be best to have spec flag names that need normalization,
		// and convert all the names at load time ( probably storing the original somewhere )
		// ( ex. store the normalized names in the meta data )
		kind: lang.Normalize(kind),
	}}
}

// defers execution; so no return value.
func (b *FieldBuilder) AddField(fn FieldInfo) {
	b.fields[PatternLocals] = append(b.fields[PatternLocals], fn)
}

func (fs *Fields) writeFields(pen *Pen) (err error) {
	if kind, e := pen.findRequiredKind(fs.kind); e != nil {
		err = e
	} else if cache, e := fs.fieldSet.cache(pen); e != nil {
		err = e
	} else if e := fs.rewriteImplicitAspects(pen, kind, &cache); e != nil {
		err = e
	} else if e := fs.fieldSet.writeFieldSet(pen, kind, cache); e != nil {
		err = e
	} else if e := fs.writeDefaultTraits(pen, kind, cache); e != nil {
		err = e
	}
	if err != nil {
		err = errutil.Fmt("%w in pattern %q domain %q", err, fs.kind, pen.domain)
	}
	return
}

// rewrite object fields
func (fs *Fields) rewriteImplicitAspects(pen *Pen, kind kindInfo, cache *classCache) (err error) {
	if isObject := strings.HasSuffix(kind.fullpath(), pen.paths.kindsPath); isObject {
		for i := range fs.fields[PatternLocals] {
			field := &fs.fields[PatternLocals][i]
			// rewrite Bool: "something" to an affinity with the opposite "not something" available.
			// i originally wanted to limit or force these into the format "is something"
			// but that screws with grok, sentences would have to be: "the noun is is something"
			if field.Affinity == affine.Bool && len(field.Class) == 0 {
				// default trait is the unset version
				defaultTrait := lang.Join([]string{"not", field.Name})
				traits := []string{defaultTrait, field.Name}
				// rewrite bool fields as implicit aspects
				aspect := lang.Join([]string{field.Name, "aspect"})
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
func (fs *Fields) writeDefaultTraits(pen *Pen, kind kindInfo, cache classCache) (err error) {
	if isObject := strings.HasSuffix(kind.fullpath(), pen.paths.kindsPath); isObject {
		for i := range fs.fields[PatternLocals] {
			field := &fs.fields[PatternLocals][i]
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

// indexed by class name
// future: wrap up the in progress field set with a "promise"
// to avoid looking up the same info repeatedly
type classCache map[string]kindInfo

func (p *classCache) store(name string, cls kindInfo) {
	if *p == nil {
		*p = make(classCache)
	}
	(*p)[name] = cls
}

// for patterns, waits to create the pattern after all fields are known
// that ensures that "extend pattern" (to add locals) happens after define pattern (for parameters and locals)
func (fs *fieldSet) cache(pen *Pen) (ret classCache, err error) {
	for ft, fields := range fs.fields {
		for _, field := range fields {
			if e := field.validate(FieldType(ft)); e != nil {
				err = e
			} else if clsName := field.getClass(); len(clsName) > 0 {
				if cls, e := pen.findOptionalKind(clsName); e != nil {
					err = e
				} else {
					ret.store(clsName, cls)
				}
			}
			if err != nil {
				err = errutil.Fmt("%w trying to find field %q", err, field.Name)
				break
			}
		}
	}
	return
}

func (fs *fieldSet) writeFieldSet(pen *Pen, kid kindInfo, cache classCache) (err error) {
	type fn func(kid, cls kindInfo, field string, aff affine.Affinity) error
	var out = [3]fn{
		pen.addParameter, // PatternParameters
		pen.addResult,
		pen.addField,
	}

Out:
	for ft, fields := range fs.fields {
		call := out[ft]
		for _, field := range fields {
			cls := cache[field.getClass()]
			e := call(kid, cls, field.Name, field.Affinity)
			if e := eatDuplicates(pen.warn, e); e != nil {
				err = e
			} else if field.Init != nil {
				e := pen.addDefaultValue(kid, field.Name, field.Init)
				if e := eatDuplicates(pen.warn, e); e != nil {
					err = e
				}
			}
			if err != nil {
				err = errutil.Fmt("%w trying to write field %q", err, field.Name)
				break Out
			}
		}
	}
	return
}

func (fs *FieldInfo) isAspectLike() (ret bool) {
	return fs.Affinity == affine.Text && fs.Name == fs.Class
}

// shortcut: if we specify a field name for a record and no class,
// we'll expect the class to be the name.
func (fs *FieldInfo) getClass() (ret string) {
	if cls := fs.Class; len(cls) > 0 {
		ret = cls
	} else if isRecordAffinity(fs.Affinity) {
		ret = fs.Name
	}
	return
}

func (fs *FieldInfo) validate(ft FieldType) (err error) {
	if len(fs.Name) == 0 {
		err = errutil.New("missing name")
	} else if len(fs.Affinity) == 0 {
		err = errutil.New("missing affinity")
	}
	return
}

// if there is a class specified, only certain affinities are allowed.
func isRecordAffinity(a affine.Affinity) (okay bool) {
	switch a {
	case affine.Record, affine.RecordList:
		okay = true
	}
	return
}
