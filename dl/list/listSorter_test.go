package list_test

import (
	"reflect"
	"testing"

	"git.sr.ht/~ionous/tapestry/dl/list"
)

// test that the sorting object's options work as expected.
func TestInnerSort(t *testing.T) {
	str := []string{"Bandana", "pair", "apple"}
	//
	list.Sorter{}.SortStrings(str)
	if !reflect.DeepEqual(str, []string{"apple", "Bandana", "pair"}) {
		t.Fatal("failed default strings", str)
	}
	//
	list.Sorter{Descending: true}.SortStrings(str)
	if !reflect.DeepEqual(str, []string{"pair", "Bandana", "apple"}) {
		t.Fatal("failed descending", str)
	}
	//
	list.Sorter{CaseSensitive: true}.SortStrings(str)
	if !reflect.DeepEqual(str, []string{"Bandana", "apple", "pair"}) {
		t.Fatal("failed case sensitive", str)
	}
	//
	list.Sorter{CaseSensitive: true}.SortStrings(str)
	if !reflect.DeepEqual(str, []string{"Bandana", "apple", "pair"}) {
		t.Fatal("failed case sensitive", str)
	}
	//
	num := []float64{16, 23, 4, 42, 8, 15}
	list.Sorter{}.SortFloats(num)
	if !reflect.DeepEqual(num, []float64{4, 8, 15, 16, 23, 42}) {
		t.Fatal("failed default nums", num)
	}
	//
	list.Sorter{Descending: true}.SortFloats(num)
	if !reflect.DeepEqual(num, []float64{42, 23, 16, 15, 8, 4}) {
		t.Fatal("failed descending nums", num)
	}
}
