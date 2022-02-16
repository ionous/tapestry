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
	} else if spec.OpenStrings {
		ret = str
	}
	return
}

// fix? it's a bit of a cheat that there's no "SetCompactValue"
// we peek at the value to see if its a key and switch on processing
// we currently store the $KEY in the in-memory values.
func (n Enum) SetValue(kv interface{}) (okay bool) {
	if str, ok := kv.(string); ok {
		spec := n.Compose()
		if len(str) > 0 && str[0] == '$' {
			// checks open strings in case we're using a $ in a keyd value
			// which we should probably handle in a more clever way
			// fix: i think OpenStrings should probably never use keys
			// and the values should just be hints ....
			if _, i := spec.IndexOfChoice(str); i >= 0 || spec.OpenStrings {
				*(n.str) = str
				okay = true
			}
		} else {
			if k, i := spec.IndexOfValue(str); i >= 0 {
				*(n.str) = k
				okay = true
			} else if spec.OpenStrings {
				*(n.str) = str
				okay = true
			}
		}
	}
	return
}
