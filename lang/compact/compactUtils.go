package compact

import "strings"

// type conversion: convert a slice containing float values into a float slice.
// returns false if any member of the slice isnt a float64
func SliceFloats(slice []any) ([]float64, bool) {
	return condition[float64](slice)
}

// type conversion: convert a slice containing string values into a string slice.
// returns false if any member of the slice isnt a string
func SliceStrings(slice []any) ([]string, bool) {
	return condition[string](slice)
}

// type conversion: convert a slice containing bool values into a bool slice.
// returns false if any member of the slice isnt a bool
func SliceBools(slice []any) ([]bool, bool) {
	return condition[bool](slice)
}

// type conversion: convert a slice of interfaces
// to a slice of strings joined with newlines.
func JoinLines(slice []any) (ret string, okay bool) {
	okay = true // provisionally
	var b strings.Builder
	for i, el := range slice {
		if str, ok := el.(string); !ok {
			okay = false
			break
		} else {
			if i > 0 {
				b.WriteRune('\n')
			}
			b.WriteString(str)
		}
	}
	if okay {
		ret = b.String()
	}
	return
}

func condition[V any](slice []any) (ret []V, okay bool) {
	okay = true // provisionally
	for i, el := range slice {
		if v, ok := el.(V); !ok {
			okay = false
			break
		} else {
			if i == 0 {
				ret = make([]V, len(slice))
			}
			ret[i] = v
		}
	}
	return
}
