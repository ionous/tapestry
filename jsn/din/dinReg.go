package din

import (
	"git.sr.ht/~ionous/iffy"
	"git.sr.ht/~ionous/iffy/dl/composer"
	"git.sr.ht/~ionous/iffy/ephemera/story"
)

var reg composer.Registry

func makeRegistry() composer.Registry {
	if reg == nil {
		for _, slats := range iffy.AllSlats {
			reg.RegisterTypes(slats)
		}
		reg.RegisterTypes(story.Slats)
	}
	return reg
}
