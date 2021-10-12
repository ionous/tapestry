package cin

import (
	"testing"

	"github.com/kr/pretty"
)

func TestSigReader(t *testing.T) {
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
	return append([]string{x.cmd}, x.params...)
}
