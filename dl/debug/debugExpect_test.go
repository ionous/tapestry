package debug

import "testing"

func TestCompareLines(t *testing.T) {
	match := []string{"a", "b", "c"}
	tests := []struct {
		okay         bool
		match, input []string
	}{
		{true, match, []string{"a", "b", "c"}},
		{false, match, []string{"x", "a", "b", "c"}},
		{false, match, []string{"b", "c"}},
		{false, match, []string{""}},
		{true, []string{""}, []string{""}},
	}
	for i, test := range tests {
		e := compareLines(test.match, test.input)
		if test.okay != (e == nil) {
			t.Fatal(i, test, e)
		}
	}
}
