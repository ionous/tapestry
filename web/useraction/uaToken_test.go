package useraction

import "testing"

func TestToken(t *testing.T) {
	from := MakeToken()
	to := ReadToken(from.String())
	if !to.Valid() || to != from {
		t.Fatalf("failed")
	}
}
