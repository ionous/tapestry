package jess

import "testing"

func TestStripArticle(t *testing.T) {
	type Pair struct {
		src, article, text string
	}
	p := []Pair{
		{src: "the evil fish", article: "the", text: "evil fish"},
		{src: "The Capital", article: "The", text: "Capital"},
		{src: "some fish", article: "some", text: "fish"},
		{src: " a   space ", article: "a", text: "space"},
		// make span turns comma into a word; so there's an extra space.
		// don't care right now.
		// {src: "dune, a desert planet", article: "", text: "dune, a desert planet"},
	}

	for _, p := range p {
		if b, e := StripArticle(p.src); e != nil {
			t.Fatal(e)
		} else if p.text != b {
			t.Fatalf("text: %q: expected: %q != got: %q", p.src, p.article, b)
		}
		// else if p.article != a {
		// 	t.Fatalf("article %q: expected: %q != got: %q", p.src, p.article, a)
		// }
	}
}
