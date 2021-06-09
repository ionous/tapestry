package decode

import "git.sr.ht/~ionous/iffy/dl/composer"

type swapType interface {
	composer.Composer
	Choices() (nameToType map[string]interface{})
}

// translate a choice, typically a $TOKEN, to a value.
// note: go-code doesnt currently have a way to find a string's label.
func FindChoice(op composer.Composer, choice string) (ret string, found bool) {
	spec := op.Compose()
	if s, i := spec.IndexOfChoice(choice); i >= 0 {
		ret = s
		found = true
	} else if spec.OpenStrings {
		ret = choice
	}
	return
}
