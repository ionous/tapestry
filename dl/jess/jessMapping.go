package jess

import (
	"git.sr.ht/~ionous/tapestry/support/match"
	"git.sr.ht/~ionous/tapestry/weave"
)

// ----
func (op *DirectionOfLinking) Match(q Query, input *InputState) (okay bool) {
	if next := *input; //
	op.Direction.Match(q, &next) &&
		op.matchFromOf(&next) &&
		op.Linking.Match(q, &next) {
		*input, okay = next, true
	}
	return
}

func (op *DirectionOfLinking) buildLink(ctx *Context) (ret *jessLink, err error) {
	if n, e := op.Linking.BuildNoun(ctx, NounProperties{}); e != nil {
		err = e
	} else {
		// direction is already normalized...
		ret = makeLink(*n, op.Direction.Text)
	}
	return
}

func (op *DirectionOfLinking) matchFromOf(input *InputState) (okay bool) {
	if m, width := fromOf.FindPrefix(input.Words()); m != nil {
		op.FromOf.Matched, *input, okay = input.Cut(width), input.Skip(width), true
	}
	return
}

var fromOf = match.PanicSpans("from", "of")

// ----
func (op *Direction) Match(q Query, input *InputState) (okay bool) {
	// options:
	// 1. look at the fields of the compass
	// 2. look at the noun instances of kind directions
	if m, width := q.FindNoun(input.Words(), Directions); width > 0 {
		op.Text = m // holds the normalized name
		*input, okay = input.Skip(width), true
	}
	return
}

// ----
func (op *Linking) Match(q Query, input *InputState) (okay bool) {
	if next := *input;        //
	op.matchNowhere(&next) || // tbd. maybe this is better than context flags? i dunno.
		Optional(q, &next, &op.KindCalled) ||
		Optional(q, &next, &op.Noun) ||
		Optional(q, &next, &op.Name) {
		*input, okay = next, true
	}
	return
}

func (op *Linking) matchNowhere(input *InputState) (okay bool) {
	if width := input.MatchWord(keywords.Nowhere); width > 0 {
		op.Nowhere = true
		*input, okay = input.Skip(width), true
	}
	return
}

// generate a room or door; an object if there's not enough information to know; or nil for nowhere.
func (op *Linking) BuildNoun(ctx *Context, props NounProperties) (ret *DesiredNoun, err error) {
	if !op.Nowhere {
		if els, e := buildNounsFrom(ctx, props, ref(op.KindCalled), ref(op.Noun), ref(op.Name)); e != nil {
			err = e
		} else {
			a := els[0]
			ret = &a
		}
	}
	return
}

// helper since we know there's linking doesnt support counted nouns, but does support nowhere;
// BuildNouns will always return a list of one or none.
func (op *Linking) GenerateNoun(ctx *Context, props NounProperties) (ret string, err error) {
	if n, e := op.BuildNoun(ctx, props); e != nil {
		err = e
	} else if n != nil {
		if e := ctx.PostProcess(weave.ValuePhase, func() error {
			return n.writeNounValues(ctx)
		}); e != nil {
			err = e
		} else {
			ret = n.Noun
		}
	}
	return
}

// ----
func (op *AdditionalLinks) Match(q Query, input *InputState) (okay bool) {
	if next := *input; //
	op.CommaAnd.Match(q, &next) &&
		op.Linking.Match(q, &next) {
		Optional(q, &next, &op.AdditionalLinks)
		*input, okay = next, true
	}
	return
}

// links iterator
type LinkIt struct {
	next  *Linking
	queue *AdditionalLinks
}

func IterateLinks(first *Linking, queue *AdditionalLinks) LinkIt {
	return LinkIt{first, queue}
}

func (it LinkIt) HasNext() bool {
	return it.next != nil
}

func (it *LinkIt) GetNext() (ret Linking) {
	ret = *it.next
	if deq := it.queue; deq == nil {
		it.next = nil
	} else {
		it.next = &deq.Linking
		it.queue = deq.AdditionalLinks
	}
	return
}

// ----
func (op *AdditionalDirections) Match(q Query, input *InputState) (okay bool) {
	if next := *input; //
	op.CommaAnd.Match(q, &next) &&
		op.DirectionOfLinking.Match(q, &next) {
		Optional(q, &next, &op.AdditionalDirections)
		*input, okay = next, true
	}
	return
}

// additional directions iterator
type DirectIt struct {
	next  *DirectionOfLinking
	queue *AdditionalDirections
}

func IterateDirections(first *DirectionOfLinking, queue *AdditionalDirections) DirectIt {
	return DirectIt{first, queue}
}

func (it DirectIt) HasNext() bool {
	return it.next != nil
}

func (it *DirectIt) GetNext() (ret DirectionOfLinking) {
	ret = *it.next
	if deq := it.queue; deq == nil {
		it.next = nil
	} else {
		it.next = &deq.DirectionOfLinking
		it.queue = deq.AdditionalDirections
	}
	return
}
