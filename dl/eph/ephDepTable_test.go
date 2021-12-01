package eph

import (
	"math/rand"
	"testing"

	"github.com/kr/pretty"
)

// sorting resolved dependencies should work as expected.
func TestTableSort(t *testing.T) {
	{
		og, tx := sorter("a", "b", "abc")
		if diff := pretty.Diff(tx, []string{"a", "b", "c"}); len(diff) > 0 {
			t.Log(pretty.Sprint(og, "->", tx))
			t.Fatal(diff)
		}
	}
	{
		og, tx := sorter("g", "gec", "gd", "ge", "gecdb", "gecdba")
		if diff := pretty.Diff(tx, []string{"g", "d", "e", "c", "b", "a"}); len(diff) > 0 {
			t.Log(pretty.Sprint(og, "->", tx))
			t.Fatal(diff)
		}
	}
	{
		og, tx := sorter("g", "gdec", "gd", "gde", "gdecb", "gdecba")
		if diff := pretty.Diff(tx, []string{"g", "d", "e", "c", "b", "a"}); len(diff) > 0 {
			t.Log(pretty.Sprint(og, "->", tx))
			t.Fatal(diff)
		}
	}
	{
		og, tx := sorter("t", "tu", "tuv", "a", "ab", "abc")
		if diff := pretty.Diff(tx, []string{"a", "b", "c", "t", "u", "v"}); len(diff) > 0 {
			t.Log(pretty.Sprint(og, "->", tx))
			t.Fatal(diff)
		}
	}
	{
		og, tx := sorter("a", "ab", "ac")
		if diff := pretty.Diff(tx, []string{"a", "b", "c"}); len(diff) > 0 {
			t.Log(pretty.Sprint(og, "->", tx))
			t.Fatal(diff)
		}
	}
}

// single letter domains, root on the left, name on the right.
// ie. "agc" means "a" domain with ancestors "g" and "c"
func sorter(strs ...string) ([]string, []string) {
	ds := make(DependencyTable, len(strs))
	for i, str := range strs {
		row := make([]Dependency, len(str))
		for j, el := range str {
			row[j] = panicDep(el)
		}
		ds[i] = Dependencies{ancestors: row}
	}
	rand.Shuffle(len(ds), func(i, j int) {
		ds[i], ds[j] = ds[j], ds[i]
	})
	was := ds.Names()
	ds.SortTable()
	return was, ds.Names()
}

type panicDep string

func (d panicDep) Name() string                           { return string(d) }
func (d panicDep) OriginAt() string                       { return "panicDep:" + string(d) }
func (d panicDep) AddRequirement(name string)             { panic("not implemented") }
func (d panicDep) GetDependencies() (Dependencies, error) { panic("not implemented") }
func (d panicDep) Resolve() (ret Dependencies, err error) { panic("not implemented") }
