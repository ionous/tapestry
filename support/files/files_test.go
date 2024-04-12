package files

import "testing"

func TestSplit(t *testing.T) {
	n, ext := SplitExt("hello.tell")
	if !ext.Story() || ext.String() != ".tell" {
		t.Fatal("bad extension", ext)
	}
	if n != "hello" {
		t.Fatal("bad split")
	}
}
