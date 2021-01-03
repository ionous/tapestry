package express

import (
	r "reflect"
	"strings"
	"testing"
	"unicode"

	"github.com/ionous/iffy"
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/dl/debug"
	"github.com/ionous/iffy/dl/list"
	"github.com/ionous/iffy/export"
	"github.com/ionous/iffy/export/tag"
	"github.com/ionous/iffy/lang"
	"github.com/kr/pretty"
)

// until template parsing gets re-written we cant handle fluid specs ( selector messaging )
// we can do a basic test to ensure it's possible to build the function signatures from the composer.Spec(s) tho.
func TestFluid(t *testing.T) {

	if got := makeSignature((*core.PutAtField)(nil)); len(pretty.Diff(got, signature{
		"put:intoRec:atField:",
		"put:intoObj:atField:",
		"put:intoObjNamed:atField:"})) > 0 {
		t.Error(strings.Join(got, ","))
	}
	if got := makeSignature((*list.SortText)(nil)); len(pretty.Diff(got, signature{
		"sort text:ascending|descending!includeCase|ignoreCase!",
		"sort text:byField:ascending|descending!includeCase|ignoreCase!"})) > 0 {
		t.Error(strings.Join(got, ","))
	}
	if got := makeSignature((*list.SortRecords)(nil)); len(pretty.Diff(got, signature{
		"sort records:using:"})) > 0 {
		t.Error(strings.Join(got, ","))
	}
	if got := makeSignature((*debug.Log)(nil)); len(pretty.Diff(got, signature{
		"log:note|toDo|warning|fix!"})) > 0 {
		t.Error(strings.Join(got, ","))
	}
}
func makeSignature(v composer.Composer) signature {
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
			unlabeled := tags.Exists("unlabeled")
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
					name := fieldName(f, tags)
					sig = sig.addSelector(name)

				case k == r.Ptr:
					// optional. so duplicate all existing selectors
					name := fieldName(f, tags)
					sig = sig.dupSelectors(name)

				case k == r.Interface && !strings.HasSuffix(n, "Eval") && n != "Assignment":
					// assumes interfaces are all unlabeled...
					slats := implementorsOf(f.Type)
					sig = sig.mulSelectors(slats)
				}
			}

		}
		return
	})
	return sig
}
func fieldName(f *r.StructField, tags tag.StructTag) (ret string) {
	if l, ok := tags.Find("label"); ok {
		ret = l
	} else {
		ret = firstRuneLower(f.Name)
	}
	return
}

func typeName(name string, t r.Type) (ret string) {
	if len(name) > 0 {
		ret = name
	} else {
		ret = lang.Underscore(t.Name())
	}
	return
}

// return the command structs supported by the passed slot
func implementorsOf(slot r.Type) (ret []string) {
	for _, slats := range iffy.AllSlats {
		for _, slat := range slats {
			slat := r.TypeOf(slat)
			if slat.Implements(slot) {
				ret = append(ret, firstRuneLower(slat.Elem().Name()))
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

func (sig signature) mulSelectors(sel []string) signature {
	var out signature
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