package internal

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

func Pascal(n string) (ret string) {
	words := strings.Split(n, "_")
	for i, w := range words {
		el, width := utf8.DecodeRuneInString(w)
		words[i] = string(unicode.ToUpper(el)) + w[width:]
	}
	return strings.Join(words, "")
}

func Camel(n string) (ret string) {
	words := strings.Split(n, "_")
	for i, w := range words {
		if i > 0 {
			el, width := utf8.DecodeRuneInString(w)
			words[i] = string(unicode.ToUpper(el)) + w[width:]
		}
	}
	return strings.Join(words, "")
}
