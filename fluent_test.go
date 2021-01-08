package iffy_test

import (
	r "reflect"
	"strings"
	"testing"
	"unicode"

	"github.com/ionous/iffy"
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/dl/list"
	"github.com/ionous/iffy/export"
	"github.com/ionous/iffy/export/tag"
	"github.com/kr/pretty"
)

// until template parsing gets re-written we cant handle fluid specs ( selector messaging )
// we can do a basic test to ensure it's possible to build the function signatures from the composer.Spec(s) tho.
func TestFluid(t *testing.T) {
	if got := makeSig((*list.Each)(nil)); !got.equals(
		"visiting:asNum:do:",
		"visiting:asTxt:do:",
		"visiting:asRec:do:",
		"visiting:asNum:do:elseIfEmptyDo:",
		"visiting:asTxt:do:elseIfEmptyDo:",
		"visiting:asRec:do:elseIfEmptyDo:",
	) {
		t.Error(got)
	}
	// if got := makeSig((*core.ChooseAction)(nil)); !got.equals(
	// 	"if:do:",
	// 	"if:do:elseIf:", // fix: sort!?
	// 	"if:do:elseDo:",
	// ) {
	// 	t.Error(got)
	// }
	// if got := makeSig((*core.Let)(nil)); !got.equals(
	// 	"let:be:",
	// ) {
	// 	t.Error(got)
	// }
	// if got := makeSig((*core.PutAtField)(nil)); !got.equals(
	// 	"put:intoRec:atField:",
	// 	"put:intoObj:atField:",
	// 	"put:intoObjNamed:atField:",
	// ) {
	// 	t.Error(got)
	// }
	// // list:
	// if got := makeSig((*list.PutIndex)(nil)); !got.equals(
	// 	"put:intoNumList:atIndex:",
	// 	"put:intoRecList:atIndex:",
	// 	"put:intoTxtList:atIndex:",
	// ) {
	// 	t.Error(got)
	// }
	// if got := makeSig((*list.PutEdge)(nil)); !got.equals(
	// 	"put:intoNumList:atBack|atFront!",
	// 	"put:intoRecList:atBack|atFront!",
	// 	"put:intoTxtList:atBack|atFront!",
	// ) {
	// 	t.Error(got)
	// }
	// if got := makeSig((*debug.Log)(nil)); !got.equals(
	// 	"log:note|toDo|warning|fix!",
	// ) {
	// 	t.Error(got)
	// }
	// if got := makeSig((*list.Erasing)(nil)); !got.equals(
	// 	"erasing:fromNumList:atIndex:as:do:",
	// 	"erasing:fromRecList:atIndex:as:do:",
	// 	"erasing:fromTxtList:atIndex:as:do:",
	// ) {
	// 	t.Error(got)
	// }
	// if got := makeSig((*list.EraseIndex)(nil)); !got.equals(
	// 	"erase:fromNumList:atIndex:",
	// 	"erase:fromRecList:atIndex:",
	// 	"erase:fromTxtList:atIndex:",
	// ) {
	// 	t.Error(got)
	// }
	// if got := makeSig((*list.EraseEdge)(nil)); !got.equals(
	// 	"erase:atBack|atFront!",
	// ) {
	// 	t.Error(got)
	// }
	// if got := makeSig((*list.Gather)(nil)); !got.equals(
	// 	"gather:fromNumList:using:",
	// 	"gather:fromRecList:using:",
	// 	"gather:fromTxtList:using:",
	// ) {
	// 	t.Error(got)
	// }
	// if got := makeSig((*list.SortText)(nil)); !got.equals(
	// 	"sort text:ascending|descending!includeCase|ignoreCase!",
	// 	"sort text:byField:ascending|descending!includeCase|ignoreCase!",
	// ) {
	// 	t.Error(got)
	// }
	// if got := makeSig((*list.SortRecords)(nil)); !got.equals(
	// 	"sort records:using:",
	// ) {
	// 	t.Error(got)
	// }
	// // rel
	// if got := makeSig((*rel.Relate)(nil)); !got.equals(
	// 	"relate obj:toObj:via:",
	// ) {
	// 	t.Error(got)
	// }
	// if got := makeSig((*rel.Reparent)(nil)); !got.equals(
	// 	"reparent childObj:toParentObj:",
	// ) {
	// 	t.Error(got)
	// }
	//

}

func (sig signature) equals(bs ...string) bool {
	return len(pretty.Diff(sig, signature(bs))) == 0
}

func (sig signature) String() string {
	return pretty.Sprint(sig)
}

func makeSig(v composer.Composer) signature {
	rtype := r.TypeOf(v).Elem()
	spec := v.Compose()
	fluid := spec.Fluent
	short := shortName(rtype, fluid.Name)

	sig := signature{short}
	// see also: parseSpec
	var cnt int
	export.WalkProperties(rtype, func(f *r.StructField, path []int) (done bool) {
		tags := tag.ReadTag(f.Tag)
		if _, ok := tags.Find("internal"); !ok {
			cnt++
			// write the selector:
			var label string
			var unlabeled bool
			if selector, ok := tags.Find("selector"); ok && len(selector) > 0 {
				label = selector
			} else {
				label = firstRuneLower(f.Name)
				unlabeled = ok
			}
			if cnt == 1 {
				var sep string
				if unlabeled {
					sep = ":"
				} else {
					sep = " "
				}
				sig[0] += sep
			}

			if cnt > 1 || !unlabeled {
				// very specific check for non-optional flags...
				if x, ok := r.Zero(r.PtrTo(f.Type)).Interface().(composer.Composer); ok {
					if spec := x.Compose(); spec.UsesStr() {
						if cs := spec.Strings; len(cs) > 0 {
							sig = sig.addFlags(cs)
							return // EARLY RETURN
						}
					}
				}
				switch n, k := f.Type.Name(), f.Type.Kind(); {
				default:
					//  write camel "fieldName:"
					sig = sig.addSelector(label)

				case k == r.Ptr:
					// optional. so duplicate all existing selectors
					sig = sig.dupSelectors(label)

				case k == r.Interface && !strings.HasSuffix(n, "Eval") && n != "Assignment":
					// assumes interfaces are all unlabeled...
					if slats := implementorsOf(f.Type); len(slats) == 0 {
						panic("no slats found") // mul will return nil
					} else {
						sig = sig.mulSelectors(slats, tags.Exists("optional"))
					}
				}
			}

		}
		return
	})
	return sig
}

// return the command structs supported by the passed slot
func implementorsOf(slot r.Type) (ret []string) {
	for _, slats := range iffy.AllSlats {
		for _, slat := range slats {
			rtype := r.TypeOf(slat)
			if rtype.Implements(slot) {
				name := composer.SpecName(slat)
				if spec := slat.Compose(); spec.Fluent != nil && len(spec.Fluent.Name) > 0 {
					name = spec.Fluent.Name
				}
				ret = append(ret, camelize(name))
			}
		}
	}
	return
}

type signature []string

func (sig signature) addSelector(sel string) signature {
	for i, cnt := 0, len(sig); i < cnt; i++ {
		sig[i] = sig[i] + sel + ":"
	}
	return sig
}

// to avoid an explosion of selectors for flags, we consider them specially.
// the signature of flags may differ from how they are specified in use, tbd.
func (sig signature) addFlags(cs []string) signature {
	var b strings.Builder
	for i, c := range cs {
		writeCamel(&b, c)
		if (i + 1) < len(cs) {
			b.WriteRune('|')
		}
	}
	b.WriteRune('!')
	sel := b.String()
	for i, cnt := 0, len(sig); i < cnt; i++ {
		sig[i] = sig[i] + sel
	}
	return sig
}

func (sig signature) mulSelectors(sel []string, optional bool) signature {
	var out signature
	if optional {
		out = sig
	}
	for i, cnt := 0, len(sig); i < cnt; i++ {
		for _, sel := range sel {
			out = append(out, sig[i]+sel+":")
		}
	}
	return out
}

func (sig signature) dupSelectors(sel string) signature {
	out := sig
	for i, cnt := 0, len(sig); i < cnt; i++ {
		out = append(out, sig[i]+sel+":")
	}
	return out
}

func shortName(t r.Type, name string) (ret string) {
	if len(name) > 0 {
		ret = name
	} else {
		ret = firstSyl(t.Name())
	}
	return
}

func firstRuneLower(s string) (ret string) {
	rs := []rune(s)
	rs[0] = unicode.ToLower(rs[0])
	return string(rs)
}

func firstSyl(s string) (ret string) {
	var b strings.Builder
	for i, u := range s {
		if i == 0 {
			b.WriteRune(unicode.ToLower(u))
		} else if !unicode.IsUpper(u) {
			b.WriteRune(u)
		} else {
			break
		}
	}
	return b.String()
}

func camelize(s string) string {
	var b strings.Builder
	writeCamel(&b, s)
	return b.String()
}

func writeCamel(b *strings.Builder, s string) {
	nextUpper := false
	for i, u := range s {
		if i == 0 {
			b.WriteRune(unicode.ToLower(u))
		} else if u == '_' || u == ' ' {
			nextUpper = true
		} else if nextUpper {
			b.WriteRune(unicode.ToUpper(u))
			nextUpper = false
		} else {
			b.WriteRune(u)
		}
	}
}