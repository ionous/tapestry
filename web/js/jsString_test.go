package js_test

import (
	"testing"

	"git.sr.ht/~ionous/tapestry/web/js"
)

func TestString(t *testing.T) {
	var out js.Builder
	out.X(`a
b`)

	if val := out.String(); val != `a\nb` {
		t.Fatal(val)
	}
}
