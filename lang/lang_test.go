package lang

import (
	"testing"
)

func TestStripArticle(t *testing.T) {
	type Pair struct {
		src, article, text string
	}
	p := []Pair{
		{src: "the evil fish", article: "the", text: "evil fish"},
		{src: "The Capital", article: "The", text: "Capital"},
		{src: "some fish", article: "some", text: "fish"},
		{src: " a   space ", article: "a", text: "space"},
		{src: "dune, a desert planet", article: "", text: "dune, a desert planet"},
	}

	for _, p := range p {
		a, b := SliceArticle(p.src)
		if p.text != b {
			t.Fatalf("text: '%s'", p.src)
		}
		if p.article != a {
			t.Fatalf("text: '%s'", p.src)
		}
	}
	if StripArticle("the article") != "article" {
		t.Fatal("article")
	}
}

// ensure package.inflect works as expected for a few cases....
// fix: they dont really work they way id expect.
func TestInflect(t *testing.T) {
	p := []struct {
		test, want string
		format     func(s string) string
	}{
		{
			"animals",
			"animals",
			Pluralize,
		},
		{
			"people",
			"people",
			Pluralize,
		},
		{
			"things",
			"things",
			Pluralize,
		},
		{
			"boop",
			"Boop",
			Capitalize,
		},
		{
			"BOOP",
			"Boop",
			Capitalize,
		}, {
			"another day. at SEA.... oh my.",
			"Another Day. At S E A.... Oh My.",
			Titlecase,
		}, {
			"another day. at SEA.... oh my.",
			"Another day. At sea.... Oh my.",
			SentenceCase,
		},
	}
	for i, p := range p {
		got := p.format(p.test)
		if got != p.want {
			t.Fatalf("test %v failed, got %q", i, got)
		}
	}
}

// go test --run TestVowels
func TestVowels(t *testing.T) {
	if !xStartsWithVowel("evil fish") {
		t.Fatal("error")
	}
}

func TestNormalize(t *testing.T) {
	if test, expect := Normalize("  hello  "), "hello"; test != expect {
		t.Error(test)
	}
	if test, expect := Normalize("Hello"), "hello"; test != expect {
		t.Error(test)
	}
	if test, expect := Normalize("Hello world"), "hello world"; test != expect {
		t.Error(test)
	}
	if test, expect := Normalize("Hello"), "hello"; test != expect {
		t.Error(test)
	}
	if test, expect := Normalize("hello_there"), "hello there"; test != expect {
		t.Error(test)
	}
	if test, expect := Normalize("hello-there"), "hello-there"; test != expect {
		t.Error(test)
	}
	if test, expect := Normalize("can't help \"this\""), "cant help this"; test != expect {
		t.Error(test)
	}
}

func TestSeparate(t *testing.T) {
	if test, expect := MixedCaseToSpaces("HelloThere"), "hello there"; test != expect {
		t.Error(test)
	} else if test, expect := MixedCaseToSpaces("helloThere"), "hello there"; test != expect {
		t.Error(test)
	}
}
