package gomake

import (
	"strings"
	"unicode"
)

// when generating some kinds of simple types...
// replace the specified typename with specified primitive.
// types that map to numbers, etc. are added as unbox automatically.
var unbox = map[string]string{"text": "string", "bool": "bool"}

// underscore_name to PascalCase
func pascal(s string) string {
	var out strings.Builder
	for _, p := range strings.Split(strings.ToLower(s), "_") {
		for i, q := range p {
			out.WriteRune(unicode.ToUpper(q))
			out.WriteString(p[i+1:])
			break
		}
	}
	return out.String()
}

// does the passed string list include the passed string?
func includes(strs []string, str string) (ret bool) {
	for _, el := range strs {
		if el == str {
			ret = true
			break
		}
	}
	return
}
