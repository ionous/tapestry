package cin

import (
	"testing"

	"github.com/kr/pretty"
)

func TestSigReader(t *testing.T) {
	// read the string on the left, do we get the parts on the right?
	for _, el := range [][]string{
		{"Story:", "story", ""},
		{"TestRule:hook:", "test_rule", "", "hook"},
		{"TestStatement testName:test:", "test_statement", "test_name", "test"},
		{"Always", "always"},
	} {
		got := readSig(el[0])
		if diff := pretty.Diff(got, el[1:]); len(diff) > 0 {
			t.Fatal(got, diff)
		}
	}
}

func readSig(s string) []string {
	var x sigReader
	x.readSig(s)
	out := []string{x.cmd}
	for _, p := range x.params {
		out = append(out, p.String())
	}
	return out
}
