package qna

import (
	"strings"
	"testing"
)

func TestSaveRandomizer(t *testing.T) {
	m := make(mockSave)
	rand := SeedRandomizer(1, 2)
	rand.Random(0, 101)
	rand.writeRandomizeer(m.save)

	// "golden image"
	if len(m) != 2 {
		t.Fatalf("expected two keys %#v", m)
	} else {
		const (
			sum  = 0xa3d02893775b61d9
			last = 77
		)
		rawSeed, rawLast := m["-$randomizer-seed"], m["-$randomizer-last"]
		if gotSum := hash(string(rawSeed.([]byte))); gotSum != sum {
			t.Fatalf("got 0x%x", gotSum)
		} else if gotLast := rawLast.(int64); gotLast != last {
			t.Fatalf("got %d", gotLast)
		}
	}
}

type mockSave map[string]any

func (m mockSave) save(domain, noun, field string, value any) (_ error) {
	key := strings.Join([]string{domain, noun, field}, "-")
	m[key] = value
	return
}
