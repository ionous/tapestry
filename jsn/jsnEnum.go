package jsn

import "git.sr.ht/~ionous/iffy/dl/composer"

// Enum wraps str-like values used by the ifspec code generator.
// it alleviates some redundant code generation.
type Enum struct {
	composer.Composer
	str *string
}

func MakeEnum(op composer.Composer, str *string) Enum {
	return Enum{op, str}
}

func (n Enum) SetEnum(kv string) bool {
	spec := n.Compose()
	if k, i := spec.IndexOfValue(kv); i >= 0 {
		*(n.str) = k
	} else {
		*(n.str) = kv
	}
	return true // fix: eventually some error handling
}

func (n Enum) GetEnum() (retKey string, retVal string) {
	spec, str := n.Compose(), *n.str
	if v, i := spec.IndexOfChoice(str); i >= 0 {
		retKey, retVal = str, v
	} else {
		retVal = str
	}
	return
}
