package rules

import (
	"testing"
)

func TestSuffixing(t *testing.T) {
	const something = "something"
	if n, s := findSuffix("something then continue"); n != something {
		t.Fatal(n)
	} else if s != Continues {
		t.Fatal(s)
	}
	if n, s := findSuffix("something then jump"); n != something {
		t.Fatal(n)
	} else if s != Jumps {
		t.Fatal(s)
	}
	if n, s := findSuffix("something then stop"); n != something {
		t.Fatal(n)
	} else if s != Stops {
		t.Fatal(s)
	}
}
