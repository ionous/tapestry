package math

import (
	"strings"
)

// Flags used compare values
type CompareType int

//go:generate stringer -type=CompareType
const (
	Compare_EqualTo CompareType = 1 << iota
	Compare_GreaterThan
	Compare_LessThan
)

// compare two floats, within some small tolerance
func (cmp CompareType) CompareFloat(a, b, epsilon float64) (ret bool) {
	return cmp.diff(compareFloats(a, b, epsilon))
}

// compare two integers
func (cmp CompareType) CompareInt(a, b int) bool {
	d := compareInt(a, b)
	return cmp.diff(d)
}

// compare two bools
func (cmp CompareType) CompareBool(a, b bool) (ret bool) {
	d := compareBool(a, b)
	return cmp.diff(d)
}

// compare two strings
func (cmp CompareType) CompareString(a, b string) bool {
	d := compareStrings(a, b)
	return cmp.diff(d)
}

func (cmp CompareType) diff(d int) (ret bool) {
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

var compareStrings = strings.Compare

func compareInt(a, b int) (ret int) {
	return a - b
}

func compareFloats(a, b, epsilon float64) (ret int) {
	switch d := a - b; {
	case d < -epsilon:
		ret = -1
	case d > epsilon:
		ret = 1
	default:
		ret = 0
	}
	return
}

func compareBool(a, b bool) (ret int) {
	switch {
	case a == b:
		ret = 0
	case a:
		ret = -1
	default:
		ret = -1
	}
	return
}
