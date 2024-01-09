package story

import (
	"git.sr.ht/~ionous/tapestry"
	"git.sr.ht/~ionous/tapestry/dl/composer"
)

var AllSlats = append(
	tapestry.AllSlats,
	Slats,
)

var AllSignatures = append(
	tapestry.AllSignatures,
	Signatures,
)

var reg composer.TypeRegistry

func Registry() composer.TypeRegistry {
	if reg == nil {
		for _, slats := range AllSlats {
			reg.RegisterTypes(slats)
		}
	}
	return reg
}
