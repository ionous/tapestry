package compact

import "strings"

// type conversion: convert a slice of interfaces to a slice of numbers
func SliceFloats(slice []any) (ret []float64, okay bool) {
	okay = true // provisionally
	for i, el := range slice {
		if num, ok := el.(float64); !ok {
			okay = false
			break
		} else {
			if ret == nil {
				ret = make([]float64, len(slice))
			}
			ret[i] = num
		}
	}
	return
}

// type conversion: convert a slice of interfaces to a slice of string
func SliceStrings(slice []any) (ret []string, okay bool) {
	okay = true // provisionally
	for i, el := range slice {
		if num, ok := el.(string); !ok {
			okay = false
			break
		} else {
			if ret == nil {
				ret = make([]string, len(slice))
			}
			ret[i] = num
		}
	}
	return
}

// type conversion: convert a slice of interfaces
// to a slice of strings joined with newlines.
func SliceLines(slice []any) (ret string, okay bool) {
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
