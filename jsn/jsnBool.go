package jsn

// Bool wraps bool-like values used by the ifspec code generator.
// it alleviates some redundant code generation.
type Bool struct {
	ComposerType
	val *bool
}

func MakeBool(op ComposerType, val *bool) Bool {
	return Bool{op, val}
}

func (n Bool) SetBool(v bool) {
	*n.val = v
}

func (n Bool) GetBool() bool {
	return *n.val
}
