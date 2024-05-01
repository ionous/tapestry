package dot

import (
	"strings"
)

// contains Index and Fields
type Path []Dotted

func (dl Path) String() string {
	var b strings.Builder
	for _, el := range dl {
		el.writeTo(&b)
	}
	return b.String()
}
