package composer

import "testing"

type Alloneword string
type CustomName struct{}
type TwoWords struct{}

func (*Alloneword) Compose() Spec { return Spec{} }
func (*CustomName) Compose() Spec { return Spec{Name: "named"} }
func (*TwoWords) Compose() Spec   { return Spec{} }

func TestNames(t *testing.T) {
	if want, have := "alloneword", SpecName((*Alloneword)(nil)); want != have {
		t.Errorf("have %q, want %q", have, want)
	} else if want, have := "named", SpecName((*CustomName)(nil)); want != have {
		t.Errorf("have %q, want %q", have, want)
	} else if want, have := "two_words", SpecName((*TwoWords)(nil)); want != have {
		t.Errorf("have %q, want %q", have, want)
	}
}
