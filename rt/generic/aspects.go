package generic

import "git.sr.ht/~ionous/tapestry/affine"

// when kinds are created via this method,
// the traits of any aspect fields act like separate boolean fields;
// without it, only the aspect text field itself exists.
// ex. if the list of fields contains a "colour" aspect with traits "red", "blue", "green"
// then the returned kind will respond to "colour", "red", "blue", and "green";
// with NewKind() it would respond only to "colour", the r,b,g fields wouldn't exist.
// its a bit of leaky abstraction because boolean traits are used only by objects.
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
