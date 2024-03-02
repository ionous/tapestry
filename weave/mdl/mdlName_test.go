package mdl_test

import (
	"reflect"
	"testing"

	"git.sr.ht/~ionous/tapestry/weave/mdl"
)

func TestNameGeneration(t *testing.T) {
	if ns := mdl.MakeNames("empire apple"); !reflect.DeepEqual(ns, []string{
		"empire apple", "apple", "empire",
	}) {
		t.Logf("mismatched %#v", ns)
		t.Fail()
	}
	if ns := mdl.MakeNames("Clarence"); !reflect.DeepEqual(ns, []string{
		"Clarence", "clarence",
	}) {
		t.Logf("mismatched %#v", ns)
		t.Fail()
	}
	if ns := mdl.MakeNames("Big Pond"); !reflect.DeepEqual(ns, []string{
		"Big Pond", "big pond", "pond", "big",
	}) {
		t.Logf("mismatched %#v", ns)
		t.Fail()
	}
}
