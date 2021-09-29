package jsn

// Str wraps str-like values used by the ifspec code generator.
// it alleviates some redundant code generation.
type Str struct {
	ComposerType
	val *string
}

func MakeStr(op ComposerType, val *string) Str {
	return Str{op, val}
}

func (n Str) SetStr(v string) {
	*n.val = v
}

func (n Str) GetStr() string {
	return *n.val
}
