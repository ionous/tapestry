package unblock

import "testing"

func TestUnstackName(t *testing.T) {
	if n, ok := unstackName("_story_break_stack"); !ok || n != "story_break" {
		t.Fatal("failed to get unstack named of type", n, ok)
	} else if n, ok := unstackName("story_break"); ok {
		t.Fatal("incorrectly unstacked named of type", n, ok)
	}
}
