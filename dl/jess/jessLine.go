package jess

import (
	"log"

	"git.sr.ht/~ionous/tapestry/lang/typeinfo"
	"git.sr.ht/~ionous/tapestry/support/match"
)

// future: maybe a plural topics for they/them; actor/s for xe/they; or inanimate topics
type topicType int

//go:generate stringer -type=topicType
const (
	undeterminedTopic topicType = iota
	otherTopic
	pronounReference // refers to some prior line's noun ( if it can )
	nounTopic
)

// a *sentence* within a paragraph. ( not a line in a file )
// fix rename to "phrase"?
type Line struct {
	words []match.TokenValue // fix? all tokens have pos; we only really need the first.
	// the successful match; mainly for debugging; its already written itself to the database
	matched typeinfo.Instance // could store MatchedPhrase maybe.
	// things we've tried; not necessarily exclusive with a valid match
	errs []error
	// pronoun helpers
	topic     GetActualNoun // might have to be GetActualNoun unless everything uses promises
	topicType topicType
}

// can return false if another phrase has determined it should be a noun.if
func (el *Line) UsePronoun() (okay bool) {
	// because its parallel matching, this can happen multiple times.
	if el.topicType == 0 || el.topicType == pronounReference {
		el.topicType = pronounReference
		okay = true
	}
	return
}

func (el *Line) SetTopic(an GetActualNoun) {
	// a successful build should (probably) only happen from one place.
	// right now, its not always clear who's going to win
	// (indicating line data, or at least this sort) should maybe be scoped to the match.
	// if el.topicType > 0 {
	// 	log.Panicf("trying to set a topic but already have %s", el.topicType)
	// }
	el.topicType = nounTopic
	el.topic = an
}

// helper for "promise" style matching.
func (el *Line) store(res PromiseMatcher) {
	if el.matched != nil {
		log.Println("matched multiple phases")
	} else {
		el.matched = res
	}
}

func (el *Line) reject(e error) {
	// fix? maybe there's a difference b/t FailedMatch and other errors?
	// Failed indicates the "shape" is wrong
	// other errors indicates the content of the particular shape is wrong.
	el.errs = append(el.errs, e)
}

// fix: make collector speak in terms of lines?
func linesToLines(src [][]match.TokenValue) []Line {
	out := make([]Line, len(src))
	for i, el := range src {
		if len(el) == 0 {
			panic("?")
		}
		out[i] = Line{words: el}
	}
	return out
}
