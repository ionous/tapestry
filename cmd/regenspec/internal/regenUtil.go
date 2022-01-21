package regen

import "strings"

func BoolOf(k string, m map[string]interface{}) (ret bool) {
	if n, ok := m[k].(bool); ok {
		ret = n
	}
	return
}

func StringOf(k string, m map[string]interface{}) (ret string) {
	if n, ok := m[k].(string); ok {
		ret = n
	}
	return
}

func MapOf(k string, m map[string]interface{}) (ret map[string]interface{}) {
	if n, ok := m[k].(map[string]interface{}); ok {
		ret = n
	}
	return
}

func ListOf(k string, m map[string]interface{}) (ret []string) {
	if ns, ok := m[k].([]interface{}); ok {
		out := make([]string, len(ns))
		for i, n := range ns {
			out[i] = n.(string)
		}
		ret = out
	}
	return
}

func sentenceOf(k string, m map[string]interface{}) (ret string) {
	n := StringOf(k, m)
	if x := strings.TrimSuffix(strings.TrimSpace(n), "."); len(x) > 0 {
		ret = x + "."
	}
	return
}

func Tokenize(k string) string {
	return "$" + strings.ToUpper(k)
}

func Detokenize(k string) string {
	return strings.ToLower(k[1:])
}
