package generic

import "git.sr.ht/~ionous/tapestry/affine"

// a bit of leaky abstraction: aspects are used only by objects
// when the aspect list exists, requests for traits of those aspects act like boolean fields
// without it, kinds can still have aspects, but only as text field enumerations.
func NewKindWithTraits(kinds Kinds, name string, parent *Kind, fields []Field) *Kind {
	return newKind(kinds, name, parent, fields, makeAspects(kinds, fields))
}

func makeAspects(kinds Kinds, fields []Field) (ret []Aspect) {
	for _, ft := range fields {
		// fix? currently a field with the same name and type is an aspect;
		// using string "aspects" might be better...
		// as there would be fewer false positives ( ex. a field of actor called actor )
		// although, it's nice the type is consistently the most derived kind...
		// ( ie. "illumination" is more specific than "aspects" )
		if ft.Affinity == affine.Text && ft.Name == ft.Type {
			if a, e := kinds.GetKindByName(ft.Type); e == nil {
				cnt := a.NumField()
				ts := make([]string, cnt)
				for i := 0; i < cnt; i++ {
					t := a.Field(i)
					ts[i] = t.Name
				}
				ret = append(ret, Aspect{Name: a.Name(), Traits: ts})
			}
		}
	}
	return
}
