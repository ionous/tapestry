package aspects

import (
	"git.sr.ht/~ionous/tapestry/affine"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
)

func MakeAspects(kinds g.Kinds, fields []g.Field) (ret []g.Aspect) {
	for _, ft := range fields {
		// tbd? currently a field with the same name and type is an aspect;
		// using string "aspects" might be better...
		// as there would be fewer false positives ( ex. a field of actor called actor )
		// although, it's nice the type is consistently the most derived kind...
		// ( ie. "illumination" is more specific than "aspects" )
		// and some of the db queries would have to change too
		if ft.Affinity == affine.Text && ft.Name == ft.Type {
			if a, e := kinds.GetKindByName(ft.Type); e == nil {
				if g.Base(a) == kindsOf.Aspect.String() {
					cnt := a.NumField()
					ts := make([]string, cnt)
					for i := 0; i < cnt; i++ {
						t := a.Field(i)
						ts[i] = t.Name
					}
					ret = append(ret, g.Aspect{Name: a.Name(), Traits: ts})
				}
			}
		}
	}
	return
}