package ident

import (
	"strings"
	"testing"
)

func TestIds(t *testing.T) {
	testEqual(t, "derived_class", nameOf("DerivedClass"))
}

func TestIdBasics(t *testing.T) {
	titleCase := nameOf("TitleCase")
	testEqual(t, "title_case", titleCase, "TitleCase to camelCase")
	testEqual(t, "two_words", nameOf("two words"), "two words to join")
	testEqual(t, "word_dash", nameOf("word-dash"), "dashes split ids")
	testEqual(t, "apostrophes", nameOf("apostrophe's"), "apostrophes vanish")
	testEqual(t, nameOf(""), "", "empty is as empty does")
	testEqual(t, nameOf("786_abc_123_def"), nameOf("786-abc 123 def"))
}

func TestIdnameOfs(t *testing.T) {
	testEqual(t, "apples", nameOf("apples"))
	testEqual(t, "apples", nameOf("Apples"))

	testEqual(t, "apple_turnover", nameOf("apple turnover"))
	testEqual(t, "apple_turnover", nameOf("Apple Turnover"))
	testEqual(t, "apple_turnover", nameOf("Apple turnover"))
	testEqual(t, "apple_turnover", nameOf("APPLE TURNOVER"))
	testEqual(t, "apple_turnover", nameOf("apple-turnover"))

	testEqual(t, "pascal_case", nameOf("PascalCase"))
	testEqual(t, "camel_case", nameOf("camelCase"))

	testEqual(t, "something_like_this", nameOf("something-like-this"))
	testEqual(t, "allcaps", nameOf("ALLCAPS"))

	// fix? hrmm... this changed at some point...
	//testEqual(t, "wha_tabout_this", nameOf("whaTAboutThis"))
	testEqual(t, "wha_t_about_this", nameOf("whaTAboutThis"))
}

func testEqual(t *testing.T, one, two string, extra ...string) {
	if one != two {
		t.Fatal(one, two, strings.Join(extra, " "))
	}
}

// TestRecycle to ensure ids generated from ids match.
// important for gopherjs optimizations.
func TestRecycle(t *testing.T) {
	src := []string{
		"lowercase",
		"ALLCAPS",
		"PascalCase",
		"camellCase",
		"space case",
		"em-dash",
	}
	for _, src := range src {
		id := nameOf(src)
		recycledId := nameOf(id)
		testEqual(t, id, recycledId)
	}
}
