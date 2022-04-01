package gomake

import (
	"strings"
	"unicode"
)

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

// underscore_name to camelCase
func camelize(s string) string {
	var out strings.Builder
	for _, p := range strings.Split(strings.ToLower(s), "_") {
		for i, q := range p {
			if out.Len() > 0 {
				q = unicode.ToUpper(q)
			}
			out.WriteRune(q)
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
