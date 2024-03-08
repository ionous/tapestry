package jess

import (
	"errors"

	"git.sr.ht/~ionous/tapestry/weave/mdl"
)

func (op *MapLocations) Match(q Query, input *InputState) (okay bool) {
	if next := *input; //
	op.Linking.Match(q, &next) &&
		op.Are.Match(q, &next) &&
		op.DirectionOfLinking.Match(q, &next) {
		Optional(q, &next, &op.AdditionalDirections)
		*input, okay = next, true
	}
	return
}

// return an iterator that is capable of walking over the right hand side of the mapping.
func (op *MapLocations) GetOtherLocations() DirectIt {
	return IterateDirections(&op.DirectionOfLinking, op.AdditionalDirections)
}

func (op *MapLocations) Generate(rar Registrar) (err error) {
	var post []postConnect
	if e := rar.PostProcess(GenerateNouns, func(q Query) (err error) {
		if a, e := op.Linking.BuildNoun(q, rar, nil, nil); e != nil {
			err = e
		} else {
			post = append(post, makePostConnect(a))
			for it := op.GetOtherLocations(); it.HasNext(); {
				link := it.GetNext()
				if b, e := link.BuildNoun(q, rar, nil, []string{""}); e != nil {
					err = e
					break
				} else {
					post = append(post, makePostConnect(b))
				}
			}
		}
		return
	}); e != nil {
		err = e
	} else if e := rar.PostProcess(GenerateDefaultKinds, func(q Query) error {
		return generateDefaultKinds(rar, post)
	}); e != nil {
		err = e
	} else {
		err = rar.PostProcess(GenerateDefaultKinds, func(q Query) (err error) {
			src, rest := post[0], post[1:]
			if !src.roomLike {
				err = src.genDoor(rar, rest)
			} else {
				err = src.genRoom(rar, rest)
			}
			return
		})
	}
	return
}

type postConnect struct {
	*DesiredNoun
	roomLike bool // valid after generating default kinds
}

func makePostConnect(n *DesiredNoun) postConnect {
	return postConnect{DesiredNoun: n}
}

func (a *postConnect) genDoor(rar Registrar, ps []postConnect) (err error) {
	for i, cnt := 0, len(ps); i < cnt && err == nil; i++ {
		if b := ps[i]; !b.roomLike {
			err = errors.New("both sides cant be doors")
		} else {
			err = a.connectDoorToRoom(rar, b)
		}
	}
	return
}

func (a *postConnect) genRoom(rar Registrar, ps []postConnect) (err error) {
	for i, cnt := 0, len(ps); i < cnt && err == nil; i++ {
		if b := ps[i]; !b.roomLike {
			err = a.connectRoomToDoor(rar, b)
		} else {
			err = a.connectRoomToRoom(rar, b)
		}
	}
	return
}

func (a *postConnect) connectDoorToRoom(rar Registrar, b postConnect) (err error) {

	// put the door (a) into room (b)
	// set the compass of (b) to (a)
	// - `B.compass[direction] = A`
	panic("not implemented")
}

func (a *postConnect) connectRoomToDoor(rar Registrar, b postConnect) (err error) {
	panic("not implemented")
}

func (a *postConnect) connectRoomToRoom(rar Registrar, b postConnect) (err error) {
	panic("not implemented")
}

func generateDefaultKinds(rar Registrar, ps []postConnect) (err error) {
	for _, p := range ps {
		if e := p.generateDefaultKind(rar); e != nil {
			err = e
			break
		}
	}
	return
}

func (p *postConnect) generateDefaultKind(rar Registrar) (err error) {
	noun := p.Noun
	if e := rar.AddNounKind(noun, Rooms); e != nil && !errors.Is(e, mdl.Duplicate) {
		err = e
	} else if e != nil {
		p.roomLike = true
	} else if e := rar.AddNounKind(noun, Doors); e != nil && !errors.Is(e, mdl.Duplicate) {
		err = e
	}
	return
}
