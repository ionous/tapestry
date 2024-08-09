package jess

import (
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

type PromiseMatcher interface {
	PromiseMatcher() PromiseMatcher
	typeinfo.Instance
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

func After(p weaver.Phase) weaver.Phase {
	return p + 1
}

func (jc JessContext) Try(z weaver.Phase, cb func(weaver.Weaves, rt.Runtime), reject func(error)) {
	if e := jc.Schedule(z, func(w weaver.Weaves, run rt.Runtime) (_ error) {
		cb(w, run)
		return
	}); e != nil {
		reject(e)
	}
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

func (jc JessContext) SetTopic(n GetActualNoun) {
	jc.CurrentPhrase().SetTopic(n)
}

func (jc JessContext) GetTopic() (ret ActualNoun) {
	p, line := jc.p, jc.phraseIndex
	for ; line >= 0; line-- {
		if el := p.Phrases[line]; el.topicType == nounTopic {
			ret = el.topic.GetActualNoun()
		} else if el.topicType != pronounReference {
			break // pronoun references will keep searching backwards
		}
	}
	return
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
