// Code generated by "stringer -type=Containment"; DO NOT EDIT.

package locate

import "fmt"

const _Containment_name = "SupportsContainsWearsCarriesHas"

var _Containment_index = [...]uint8{0, 8, 16, 21, 28, 31}

func (i Containment) String() string {
	if i < 0 || i >= Containment(len(_Containment_index)-1) {
		return fmt.Sprintf("Containment(%d)", i)
	}
	return _Containment_name[_Containment_index[i]:_Containment_index[i+1]]
}
