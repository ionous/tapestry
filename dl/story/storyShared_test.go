package story_test

import (
	"strings"
)

func lines(s ...string) string {
	return strings.Join(s, "\n") + "\n"
}
