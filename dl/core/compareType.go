package core

type CompareType int

// Comparator generates comparison flags.
// FIX: a combo-box of enumeration options should be possible.
type Comparator interface {
	Compare() CompareType
}

func (*Equal) Compare() CompareType {
	return Compare_EqualTo
}
func (*Unequal) Compare() CompareType {
	return Compare_GreaterThan | Compare_LessThan
}
func (*GreaterThan) Compare() CompareType {
	return Compare_GreaterThan
}
func (*LessThan) Compare() CompareType {
	return Compare_LessThan
}
func (*AtLeast) Compare() CompareType {
	return Compare_GreaterThan | Compare_EqualTo
}
func (*AtMost) Compare() CompareType {
	return Compare_LessThan | Compare_EqualTo
}

//go:generate stringer -type=CompareType
const (
	Compare_EqualTo CompareType = 1 << iota
	Compare_GreaterThan
	Compare_LessThan
)
