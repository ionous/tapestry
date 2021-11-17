package eph

type Conflict struct {
	Reason ReasonForConflict
	Domain string
	Def    Definition
}

func (n *Conflict) Error() string {
	return n.Reason.String()
}

type ReasonForConflict int

//go:generate stringer -type=ReasonForConflict
const (
	Redefined ReasonForConflict = iota
	Duplicated
)

type defines map[string]Definition

type Definition struct {
	at, value string
	err       error
}
