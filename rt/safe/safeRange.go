package safe

import "fmt"

// min inclusive, max exclusive
func Range(i, min, max int) (ret int, err error) {
	if i < min {
		ret, err = min, fmt.Errorf("underflow %d < %d", i, min)
	} else if i >= max {
		ret, err = max-1, fmt.Errorf("overflow %d  >= %d", i, max)
	} else {
		ret = i
	}
	return
}
