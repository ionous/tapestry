package rt

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
)

// given a set of fields which may contain aspects generate a list of traits.
// relies on the passed kinds database having types which inherit from kindOf.Aspect
// the fields of those types are the traits.
func MakeAspects(ks Kinds, fields []Field) (ret []Aspect) {
	for _, ft := range fields {
		// tbd? currently a field with the same name and type is an aspect;
		// using string "aspects" might be better...
		// as there would be fewer false positives ( ex. a field of actor called actor )
		// although, it's nice the type is consistently the most derived kind...
		// ( ie. "illumination" is more specific than "aspects" )
		// and some of the db queries would have to change too
		if ft.Affinity == affine.Text && ft.Name == ft.Type {
			if a, e := ks.GetKindByName(ft.Type); e == nil {
				if a.Implements(kindsOf.Aspect.String()) {
					cnt := a.FieldCount()
					ts := make([]string, cnt)
					for i := 0; i < cnt; i++ {
						t := a.Field(i)
						ts[i] = t.Name
					}
					ret = append(ret, Aspect{Name: a.Name(), Traits: ts})
				}
			}
		}
	}
	return
}
