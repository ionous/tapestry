package lang

import (
	"testing"
)

func TestUnderscore(t *testing.T) {
	pairs := []string{
		"", "",
		"A", "a",
		"a", "a",
		"apples", "apples", // single words
		"Apples", "apples",
		"APPLES", "apples",
		"appleTurnover", "apple_turnover", // multi-words,
		"apple turnover", "apple_turnover",
		"Apple Turnover", "apple_turnover",
		"Apple turnover", "apple_turnover",
		"APPLE TURNOVER", "apple_turnover",
		"apple-turnover", "apple_turnover",
		"apple---turn---over", "apple_turn_over",
		"WasPascalCase", "was_pascal_case", // multi-word casing,
		"wasCamelCase", "was_camel_case",
		"something-like-this", "something_like_this",
		"something_like_that", "something_like_that",
		"some___thing__like_that", "some_thing_like_that",
		"whaTAboutThis", "wha_t_about_this", // rando,
		"lowercase", "lowercase",
		"Class", "class",
		"MyClass", "my_class",
		"MyC", "my_c",
		"PDFLoader", "pdf_loader",
		"AString", "a_string",
		"SimpleXMLParser", "simple_xml_parser",
		"vimRPCPlugin", "vim_rpc_plugin",
		"GL11Version", "gl_11_version",
		"99Bottles", "99_bottles",
		"May5", "may_5",
		"BFG9000", "bfg_9000",
		"BöseÜberraschung", "böse_überraschung",
		"Two  spaces", "two_spaces",
		"BadUTF8\xe2\xe2\xa1", "bad_utf_8",
	}
	for i, cnt := 0, len(pairs); i < cnt; i += 2 {
		test, want := pairs[i], pairs[i+1]
		if got := Underscore(test); got != want {
			t.Error("line", i/2+9, "wanted", want, "got", got)
		}
	}
}
