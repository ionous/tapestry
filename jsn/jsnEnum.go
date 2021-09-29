package jsn

import "git.sr.ht/~ionous/iffy/dl/composer"

type ComposerType interface {
	composer.Composer
	GetType() string
}

// Enum wraps str-like values used by the ifspec code generator.
// it alleviates some redundant code generation.
type Enum struct {
	ComposerType
	str *string
}

func MakeEnum(op ComposerType, str *string) Enum {
	return Enum{op, str}
}

func (n Enum) SetEnum(kv string) {
	spec := n.Compose()
	if k, i := spec.IndexOfValue(kv); i >= 0 {
		*(n.str) = k
	} else {
		*(n.str) = kv
	}
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
