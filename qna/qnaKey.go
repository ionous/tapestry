package qna

import (
	"git.sr.ht/~ionous/iffy/lang"
	"git.sr.ht/~ionous/iffy/object"
	g "git.sr.ht/~ionous/iffy/rt/generic"
)

type keyType struct {
	target, field string
}

func makeKey(target, field string) keyType {
	// FIX?
	// operations generating get field should be registering the field as a name
	// and, as best as possible, relating obj to field for property verification
	// name translation should be done there.
	// we'd have to mark up things like text evaluations ( ex. HasTrait )
	if len(field) > 0 && field[0] != '#' {
		field = lang.Breakcase(field)
	}
	return keyType{target, field}
}

// FIX: remove.
func makeKeyForEval(obj, typeName string) keyType {
	return keyType{obj, typeName}
}

// return an error saying that the value of this key is unknown
func (k *keyType) unknown() (err error) {
	if k.target == object.Value {
		err = g.UnknownObject(k.field)
	} else {
		err = g.UnknownField(k.target, k.field)
	}
	return
}

// subField should be one of the package object prefixes
func (k *keyType) dot() string {
	return k.target + "." + k.field
}
