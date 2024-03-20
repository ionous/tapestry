package jess

import (
	"errors"

	"git.sr.ht/~ionous/tapestry/weave/mdl"
)

// runs in the MappingPhase phase
func (op *MapDirections) Phase() Phase {
	return mdl.MappingPhase
}

func (op *MapDirections) Match(q Query, input *InputState) (okay bool) {
	if next := *input; //
	op.DirectionOfLinking.Match(q, &next) &&
		op.Are.Match(q, &next) &&
		(Optional(q, &next, &op.Redirect) ||
			Optional(q, &next, &op.Linking)) {
		*input, okay = next, true
	}
	return
}

func (op *MapDirections) Generate(ctx *Context) (err error) {
	if op.Linking != nil {
		err = op.simpleLink(ctx)
	} else if op.Redirect != nil {
		err = op.multiLink(ctx)
	} else {
		panic("unhandled map direction")
	}
	return
}

// uses .Linking
func (op *MapDirections) simpleLink(ctx *Context) (err error) {
	var links []jessLink
	if lhs, e := op.DirectionOfLinking.buildLink(ctx); e != nil {
		err = e
	} else if rhs, e := op.Linking.BuildNoun(ctx, nil, nil); e != nil {
		err = e
	} else {
		rhs := makeLink(rhs, "")
		links = []jessLink{lhs, rhs}
		err = ctx.PostProcess(mdl.ConnectionPhase, func() (err error) {
			if e := assignDefaultKinds(ctx, links); e != nil {
				err = e
			} else {
				err = connectPlaceToPlaces(ctx, links[1], links[:1])
			}
			return
		})
	}
	return
}

// uses .Redirect
func (op *MapDirections) multiLink(ctx *Context) (err error) {
	if lhs, e := op.DirectionOfLinking.buildLink(ctx); e != nil {
		err = e
	} else if rhs, e := op.Redirect.buildLink(ctx); e != nil {
		err = e
	} else {
		err = ctx.PostProcess(mdl.ConnectionPhase, func() (err error) {
			links := []jessLink{lhs, rhs}
			if e := assignDefaultKinds(ctx, links); e != nil {
				err = e
			} else {
				// "from l is redirect r" is
				// "from l is r" and "l is redirect r" ( aka: "redirect r is l" )
				lhs, rhs := links[0], links[1]
				from, redirect := links[0].direction, links[1].direction
				if door, e := createPrivateDoor(ctx, from, lhs, rhs); e != nil {
					err = e
				} else if len(door) == 0 {
					err = errors.New("room already has a door")
				} else if door, e := createPrivateDoor(ctx, redirect, rhs, lhs); e != nil {
					err = e
				} else if len(door) == 0 {
					err = errors.New("room already has a door")
				}
			}
			return
		})
	}
	return
}
