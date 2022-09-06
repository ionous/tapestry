package core

// Flags used compare values
type CompareType int

//go:generate stringer -type=CompareType
const (
	Compare_EqualTo CompareType = 1 << iota
	Compare_GreaterThan
	Compare_LessThan
)

// compare the passed float to zero, within some small tolerance
func (cmp CompareType) CompareFloat(d, epsilon float64) (ret bool) {
	switch {
	case d < -epsilon:
		ret = (cmp & Compare_LessThan) != 0
	case d > epsilon:
		ret = (cmp & Compare_GreaterThan) != 0
	default:
		ret = (cmp & Compare_EqualTo) != 0
	}
	return
}

// compare the passed integer to zero
func (cmp CompareType) CompareInt(d int) (ret bool) {
	switch {
	case d < 0:
		ret = (cmp & Compare_LessThan) != 0
	case d > 0:
		ret = (cmp & Compare_GreaterThan) != 0
	default:
		ret = (cmp & Compare_EqualTo) != 0
	}
	return
}
