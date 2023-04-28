package weave

import "testing"

func TestOppositePairs(t *testing.T) {
	var ps OppositePairs
	if e := ps.AddPair("east", "west"); e != nil {
		t.Fatal(e)
	} else if e := ps.AddPair("east", "west"); e != nil {
		t.Fatal(e) // add the duplicate pairing should be fine
	} else if e := ps.AddPair("west", "east"); e != nil {
		t.Fatal(e) // add the inverse pairing should be fine
	} else if e := ps.AddPair("north", "south"); e != nil {
		t.Fatal(e) // add more should be fine
	} else if e := ps.AddPair("north", "inside"); e == nil {
		t.Fatal("left conflict") // conflicting words should fail
	} else if e := ps.AddPair("outside", "south"); e == nil {
		t.Fatal("right conflict") // conflicting words should fail
	} else if p, ok := ps.FindOpposite("west"); !ok || p != "east" {
		t.Fatal(p)
	} else if p, ok := ps.FindOpposite("east"); !ok || p != "west" {
		t.Fatal(p)
	}
}
