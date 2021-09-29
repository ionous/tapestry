package jsn

// Num wraps float64-like values used by the ifspec code generator.
// it alleviates some redundant code generation.
type Num struct {
	ComposerType
	val *float64
}

func MakeNum(op ComposerType, val *float64) Num {
	return Num{op, val}
}

func (n Num) SetNum(v float64) {
	*n.val = v
}

func (n Num) GetNum() float64 {
	return *n.val
}
