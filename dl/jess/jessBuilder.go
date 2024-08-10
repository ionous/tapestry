package jess

import (
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/weave/weaver"
)

// the process of matching an english phrase creates a builder.
// after the first successful match, the associated builder is triggered.
type Builder interface {
	Build(BuildContext) error
}

type BuildContext struct {
	Query
	weaver.Weaves
	run rt.Runtime
}

func (an ActualNoun) BuildPropertyNoun(BuildContext) string {
	return an.Name
}
