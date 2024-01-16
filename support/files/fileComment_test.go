package files

import (
	r "reflect"
	"testing"
)

// verify the line splitting and stripping works as expected
func TestReadComment(t *testing.T) {
	one := readComment("# hello")
	if str, ok := one.(string); !ok || str != "hello" {
		t.Fatal(ok, str)
	}
	many := readComment("# hello\n# there\n# world")
	if !r.DeepEqual(many, []any{"hello", "there", "world"}) {
		t.Fatalf("%#v", many)
	}
}
