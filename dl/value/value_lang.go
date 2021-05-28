// Code generated by "makeops"; edit at your own risk.
package value

import (
	"git.sr.ht/~ionous/iffy/dl/composer"
)

// Text requires a user-specified string.
type Text string

func (*Text) Choices() (choices map[string]string) {
	return map[string]string{
		"$EMPTY": "empty",
	}
}

func (*Text) Compose() composer.Spec {
	return composer.Spec{
		Name:        "text",
		OpenStrings: true,
		Strings: []string{
			"empty",
		},
	}
}
