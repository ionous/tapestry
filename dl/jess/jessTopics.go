package jess

import (
	"errors"
	"fmt"

	"git.sr.ht/~ionous/tapestry/weave/weaver"
)

// has three states: contains error, contains a noun, contains either.
type Topic struct {
	noun ActualNoun
	err  error // for polling ( non-promises )
}

// doesn't fire any callbacks ( that's handled at the phrase level )
func (tc *Topic) rejectedTopic(e error) {
	if tc.err == nil { // keep the first error
		tc.err = e
	}
}

// doesn't fire any callbacks ( that's handled at the phrase level )
func (tc *Topic) acceptedTopic(an ActualNoun) {
	if !an.IsValid() {
		tc.err = errors.New("invalid noun")
	} else {
		tc.noun = an
	}
}

// returns errMissing so the scheduler will loop until the topic completes or errors.
func (tc *Topic) Resolve() (ret ActualNoun, err error) {
	if e := tc.err; e != nil {
		err = e
	} else if an := tc.noun; !an.IsValid() {
		err = fmt.Errorf("%w pronoun topic", weaver.ErrMissing)
	} else {
		ret = an
	}
	return
}

// triggered after build because that's when topics are formed
type promisedTopic struct {
	accept func(ActualNoun)
	reject func(error)
}

type promisedTopics []promisedTopic

// if the noun isn't valid, this turns into a generic rejection
func (pt promisedTopics) accept(an ActualNoun) {
	if !an.IsValid() {
		e := fmt.Errorf("there was an unknown error determining the topic")
		pt.reject(e)
	} else {
		for _, el := range pt {
			el.accept(an)
		}
	}
}

func (pt promisedTopics) reject(err error) {
	for _, el := range pt {
		el.reject(err)
	}
}
