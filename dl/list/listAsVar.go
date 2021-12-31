package list

import (
	"git.sr.ht/~ionous/tapestry/affine"
)

// ListIterator defines a variable name
type ListIterator interface {
	// ? maybe we could be the thing that genates the record?
	Name() string
	Affinity() affine.Affinity
}

func (op *AsNum) Name() string {
	return op.Var.String()
}

func (op *AsRec) Name() string {
	return op.Var.String()
}

func (op *AsTxt) Name() string {
	return op.Var.String()
}

func (op *AsNum) Affinity() affine.Affinity {
	return affine.Number
}

func (op *AsRec) Affinity() affine.Affinity {
	return affine.Record
}

func (op *AsTxt) Affinity() affine.Affinity {
	return affine.Text
}
