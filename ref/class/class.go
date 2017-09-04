package class

import (
	"github.com/ionous/iffy/ident"
	"github.com/ionous/iffy/lang"
	"github.com/ionous/iffy/ref/unique"
	r "reflect"
)

// Id returns a normalized identifier. Normalized identifiers are often used in registries to find types by name.
func Id(rtype r.Type) ident.Id {
	name := rtype.Name()
	return ident.IdOf(name)
}

// FriendlyName returns a printable string.
func FriendlyName(rtype r.Type) string {
	name := rtype.Name()
	return lang.Lowerspace(name)
}

// Parent types are designated by a struct field with the tag `if:"parent"`.
// This allows for single inheritance in pod types, similar to c-struct embedding.
// ex. type DerivedClass struct { BaseClass `if:"parent"` }
func Parent(rtype r.Type) (ret r.Type, okay bool) {
	if path, ok := unique.PathOf(rtype, "parent"); ok {
		ret, okay = rtype.FieldByIndex(path).Type, true
	}
	return
}

// IsSame tests whether the passed name refers to the passed type.
func IsSame(rtype r.Type, name string) (okay bool) {
	if id := ident.IdOf(name); Id(rtype) == id {
		okay = true
	}
	return
}

// IsCompatible tests whether the passed type can be used as the named class.
// ie. is the named type a Parent() of the passed type?
// FIX? would this be better as Type vs. Type, leaving the name -> Type lookup as part of a registry.
func IsCompatible(rtype r.Type, name string) (okay bool) {
	if id := ident.IdOf(name); Id(rtype) == id {
		okay = true
	} else {
		for i := 0; i < 250; i++ {
			if p, ok := Parent(rtype); !ok {
				break
			} else if Id(p) == id {
				okay = true
				break
			} else {
				rtype = p
			}
		}
	}
	return
}

// PropertyPath searches the passed type for a field with the passed name.
func PropertyPath(rtype r.Type, name string) (ret []int) {
	pid := ident.IdOf(name)
	fn := func(f *r.StructField, path []int) (done bool) {
		if ident.IdOf(f.Name) == pid {
			ret, done = path, true
		}
		return
	}
	unique.WalkProperties(rtype, fn)
	return
}
