package jsn

import "git.sr.ht/~ionous/tapestry/dl/composer"

// Enum wraps str-like values used by the ifspec code generator.
// it alleviates some redundant code generation.
type Enum struct {
	composer.Composer
	str *string
}

func MakeEnum(op composer.Composer, str *string) Enum {
	return Enum{op, str}
}

func (n Enum) GetValue() interface{} {
	return *n.str
}

func (n Enum) String() string {
	return *n.str
}

func (n Enum) GetCompactValue() (ret interface{}) {
	spec, str := n.Compose(), *n.str
	if v, i := spec.IndexOfChoice(str); i >= 0 {
		ret = v
	} else {
		ret = str
	}
	return
}

func (n Enum) SetValue(kv interface{}) (okay bool) {
	if str, ok := kv.(string); ok {
		spec := n.Compose()
		if k, i := spec.IndexOfValue(str); i >= 0 {
			*(n.str) = k
		} else {
			*(n.str) = str
		}
		okay = true
	}
	return
}
