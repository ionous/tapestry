package generate

import "testing"

func TestPascal(t *testing.T) {
	if a := Pascal("underscore_name"); a != "UnderscoreName" {
		t.Fatal("mismatch", a)
	}
}

func TestCamalize(t *testing.T) {
	if a := Camelize("underscore_name"); a != "underscoreName" {
		t.Fatal("mismatch", a)
	}
}
