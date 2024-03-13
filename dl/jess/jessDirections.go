package jess

import (
	"errors"
)

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

func (op *MapDirections) Generate(rar *Context) (err error) {
	if op.Linking != nil {
		err = op.simpleLink(rar)
	} else if op.Redirect != nil {
		err = op.multiLink(rar)
	} else {
		panic("unhandled map direction")
	}
	return
}

// uses .Linking
func (op *MapDirections) simpleLink(rar *Context) (err error) {
	var links []jessLink
	if e := rar.PostProcess(GenerateNouns, func(q Query) (err error) {
		if lhs, e := op.DirectionOfLinking.buildLink(q, rar); e != nil {
			err = e
		} else if rhs, e := op.Linking.BuildNoun(q, rar, nil, nil); e != nil {
			err = e
		} else {
			rhs := makeLink(rhs, "")
			links = []jessLink{lhs, rhs}
		}
		return
	}); e != nil {
		err = e
	} else if e := rar.PostProcess(GenerateDefaultKinds, func(Query) error {
		return generateDefaultKinds(rar, links)
	}); e != nil {
		err = e
	} else {
		err = rar.PostProcess(GenerateConnections, func(Query) error {
			return connectPlaceToPlaces(rar, links[1], links[:1])
		})
	}
	return
}

// uses .Redirect
func (op *MapDirections) multiLink(rar *Context) (err error) {
	var links []jessLink
	if e := rar.PostProcess(GenerateNouns, func(q Query) (err error) {
		if lhs, e := op.DirectionOfLinking.buildLink(q, rar); e != nil {
			err = e
		} else if rhs, e := op.Redirect.buildLink(q, rar); e != nil {
			err = e
		} else {
			links = []jessLink{lhs, rhs}
		}
		return
	}); e != nil {
		err = e
	} else if e := rar.PostProcess(GenerateDefaultKinds, func(Query) error {
		return generateDefaultKinds(rar, links)
	}); e != nil {
		err = e
	} else {
		// "from l is redirect r" is
		// "from l is r" and "l is redirect r" ( aka: "redirect r is l" )
		err = rar.PostProcess(GenerateConnections, func(Query) (err error) {
			lhs, rhs := links[0], links[1]
			from, redirect := links[0].direction, links[1].direction
			if door, e := createPrivateDoor(rar, from, lhs, rhs); e != nil {
				err = e
			} else if len(door) == 0 {
				err = errors.New("room already has a door")
			} else if door, e := createPrivateDoor(rar, redirect, rhs, lhs); e != nil {
				err = e
			} else if len(door) == 0 {
				err = errors.New("room already has a door")
			}
			return
		})
	}
	return
}
