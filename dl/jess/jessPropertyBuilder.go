package jess

import "git.sr.ht/~ionous/tapestry/rt"

type PropertyBuilder struct {
	Context JessContext // maybe
	noun    PropertyNoun
	pending []pendingProperty
}

func (pb *PropertyBuilder) GetKind() string {
	return pb.noun.GetKind()
}

func (pb *PropertyBuilder) addProperty(n string, v rt.Assignment) {
	pb.pending = append(pb.pending, pendingProperty{n, v})
}

type pendingProperty struct {
	fieldName string
	val       rt.Assignment
}

func (pb *PropertyBuilder) Build(out BuildContext) (err error) {
	if noun, e := pb.noun.BuildPropertyNoun(out); e != nil {
		err = e
	} else {
		for _, prop := range pb.pending {
			if e := out.AddNounValue(noun, prop.fieldName, prop.val); e != nil {
				err = e
				break
			}
		}
	}
	return
}
