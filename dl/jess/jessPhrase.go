package jess

import (
	"log"

	"git.sr.ht/~ionous/tapestry/lang/typeinfo"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/support/match"
	"git.sr.ht/~ionous/tapestry/weave/weaver"
)

// a sentence within a paragraph.
type Phrase struct {
	words []match.TokenValue // fix? all tokens have pos; we only really need the first.
	// the successful match; mainly for debugging; its already written itself to the database
	matched  typeinfo.Instance // could store MatchedPhrase maybe.
	promised PromisedMatcher
	// things we've tried; not necessarily exclusive with a valid match
	errs []error

	// pronoun helpers
	topic         Topic
	pendingTopics promisedTopics
}

// returns true if a build was initiated
func (el *Phrase) Build(jc JessContext) (okay bool) {
	if p := el.promised; p != nil {
		if useLogging(jc) {
			m := Matched(el.words)
			log.Printf("matched %s %q\n", p.TypeInfo().TypeName(), m.DebugString())
		}
		jc.Schedule(weaver.AnyPhase, func(w weaver.Weaves, run rt.Runtime) (err error) {
			b := p.GetBuilder() // builds in the "background"
			if e := b.Build(BuildContext{jc.Query, w, run}); e != nil {
				el.pendingTopics.reject(e)
				err = e
			} else {
				el.pendingTopics.accept(el.topic.noun)
			}
			return
		})
		okay = true
	}
	return
}

// helper for "promise" style matching.
func (el *Phrase) store(res PromisedMatcher) {
	if el.matched != nil || el.promised != nil {
		log.Println("matched multiple patterns")
	} else {
		el.promised = res
	}
}

func (el *Phrase) reject(e error) {
	// fix? maybe there's a difference b/t FailedMatch and other errors?
	// Failed indicates the "shape" is wrong
	// other errors indicates the content of the particular shape is wrong.
	el.errs = append(el.errs, e)
}

// fix: make collector speak in terms of phrases?
func tokensToPhrases(src [][]match.TokenValue) []Phrase {
	out := make([]Phrase, len(src))
	for i, el := range src {
		if len(el) == 0 {
			panic("?")
		}
		out[i] = Phrase{words: el}
	}
	return out
}
