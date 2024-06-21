package query

import (
	"testing"
)

func TestSaveRandomizer(t *testing.T) {
	rand := SeedRandomizer(1, 2)
	rand.Random(0, 101)
	if p, e := rand.Save(); e != nil {
		t.Fatal(e)
	} else {
		const (
			sum  = 0xa3d02893775b61d9
			last = 77
		)
		rawSeed, rawLast := p.Seed, p.LastRandom
		if gotSum := hash(string(rawSeed)); gotSum != sum {
			t.Fatalf("got 0x%x", gotSum)
		} else if gotLast := rawLast; gotLast != last {
			t.Fatalf("got %d", gotLast)
		}
	}
}
