package generate

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

// underscore_name to UnderscoreName
func Pascal(s string) string {
	var out strings.Builder
	for _, p := range strings.Split(strings.ToLower(s), "_") {
		if q, w := utf8.DecodeRuneInString(p); w > 0 {
			out.WriteRune(unicode.ToUpper(q))
			out.WriteString(p[w:])
		}
	}
	return out.String()
}

// underscore_name to underscoreName
func Camelize(s string) string {
	var out strings.Builder
	for _, p := range strings.Split(strings.ToLower(s), "_") {
		if q, w := utf8.DecodeRuneInString(p); w > 0 {
			if out.Len() > 0 {
				q = unicode.ToUpper(q)
			}
			out.WriteRune(q)
			out.WriteString(p[w:])
		}
	}
	return out.String()
}

// does the passed string list include the passed string?
func Includes(strs []string, str string) (ret bool) {
	for _, el := range strs {
		if el == str {
			ret = true
			break
		}
	}
	return
}
