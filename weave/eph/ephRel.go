package eph

type Cardinality interface {
	Cardinality()
}

// OneOne
type OneOne struct {
	Kind      string `if:"label=_,type=text"`
	OtherKind string `if:"label=to_kind,type=text"`
}

func (*OneOne) Cardinality() {}

// OneMany
type OneMany struct {
	Kind       string `if:"label=_,type=text"`
	OtherKinds string `if:"label=to_kinds,type=text"`
}

func (*OneMany) Cardinality() {}

// ManyMany
type ManyMany struct {
	Kinds      string `if:"label=_,type=text"`
	OtherKinds string `if:"label=to_kinds,type=text"`
}

func (*ManyMany) Cardinality() {}

// ManyOne
type ManyOne struct {
	Kinds     string `if:"label=_,type=text"`
	OtherKind string `if:"label=to_kind,type=text"`
}

func (*ManyOne) Cardinality() {}
