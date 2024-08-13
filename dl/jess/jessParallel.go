package jess

import (
	"errors"

	"git.sr.ht/~ionous/tapestry/lang/compact"
	"git.sr.ht/~ionous/tapestry/lang/typeinfo"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/support/match"
	"git.sr.ht/~ionous/tapestry/weave/weaver"
)

type Scheduler interface {
	// can return error if out of phase
	Schedule(weaver.Phase, func(weaver.Weaves, rt.Runtime) error) error
	SchedulePos(compact.Source, weaver.Phase, func(weaver.Weaves, rt.Runtime) error) error
}

func After(p weaver.Phase) weaver.Phase {
	return p + 1
}

// what do accept/reject function pairs generate?
// they generate builders.
type PromisedMatcher interface {
	GetBuilder() Builder
	typeinfo.Instance // helps with logging: they are all generated classes.
}

type JessContext struct {
	Query
	Scheduler
	// note: the paragraph and phrase appear in the input state as well
	// probably should just be there, and not here.
	p           *Paragraph
	phraseIndex int
	flags       int
}

// override query context to provide additional flags
// see AddContext, ClearContext
func (jc JessContext) GetContext() int {
	return jc.flags | jc.Query.GetContext()
}

func (jc JessContext) CurrentPhrase() *Phrase {
	return &jc.p.Phrases[jc.phraseIndex]
}

func (jc JessContext) Source() (ret compact.Source) {
	p, el := jc.p, jc.CurrentPhrase()
	lineOfs := el.words[0].Pos.Y
	return compact.Source{
		File:    p.File,
		Line:    lineOfs,
		Comment: "a plain-text paragraph",
	}
}

// overrides the normal find noun to map "you" to the object "self"
// fix: it'd be nice if the mapping of "you" to "self" was handled by script
// ( ex. registering the appropriate names )
func (jc JessContext) FindNoun(name []match.TokenValue, pkind *string) (string, int) {
	if len(name) == 1 {
		n := name[0] // a copy
		if n.Token == match.String && n.Hash() == keywords.You {
			n.Value = PlayerSelf         // replace in the copy
			name = []match.TokenValue{n} // leave the original data as is.
		}
	}
	return jc.Query.FindNoun(name, pkind)
}

// run the passed function now or in the future.
func (jc JessContext) Schedule(when weaver.Phase, cb func(weaver.Weaves, rt.Runtime) error) (err error) {
	pos := jc.Source()
	return jc.Scheduler.SchedulePos(pos, when, cb)
}

func (jc JessContext) Try(z weaver.Phase, cb func(weaver.Weaves, rt.Runtime), reject func(error)) {
	if e := jc.Schedule(z, func(w weaver.Weaves, run rt.Runtime) (_ error) {
		cb(w, run)
		return
	}); e != nil {
		reject(e)
	}
}

func (jc JessContext) TryTopic(accept func(ActualNoun), reject func(error)) {
	if p, i := jc.p, jc.phraseIndex-1; i < 0 {
		e := errors.New("has no preceding sentence")
		reject(e)
	} else {
		el := &p.Phrases[i]
		if e := el.topic.err; e != nil {
			reject(e)
		} else if el.topic.noun.IsValid() {
			accept(el.topic.noun)
		} else {
			el.pendingTopics = append(el.pendingTopics, promisedTopic{accept, reject})
		}
	}
}

// ideally this would be in BuildContext, or even returned from Build()
// but not not everything uses that.
func (jc JessContext) SetTopic(an ActualNoun) {
	el := jc.CurrentPhrase()
	el.topic.acceptedTopic(an)
	// break a chain of deeper callbacks by scheduling to the "any phase"
	jc.Schedule(0, func(weaver.Weaves, rt.Runtime) error {
		el.pendingTopics.accept(el.topic.noun)
		return nil
	})
}

// when a sentence doesn't use a topic
func (jc JessContext) RejectTopic(e error) {
	el := jc.CurrentPhrase()
	el.topic.rejectedTopic(e)
	jc.Schedule(0, func(weaver.Weaves, rt.Runtime) error {
		el.pendingTopics.reject(el.topic.err)
		return nil
	})
}
