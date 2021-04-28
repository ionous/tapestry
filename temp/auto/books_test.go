package books

import (
	"testing"

	"zombiezen.com/go/capnproto2/std/capnp/schema"
)

func TestBooks(t *testing.T) {
	var nodes Map // uses the schemas.DefaultRegistry by default

	n, e := nodes.Find(Squeak_TypeID)
	if e != nil {
		t.Fatal(e)
	}
	if !n.IsValid() || n.Which() != schema.Node_Which_structNode {
		t.Fatal("cannot find struct type")
	}
	x := n.StructNode()
	// where "s" is a valid struct.
	// var discriminant uint16
	// if hasDiscriminant(n) {
	// 	// byte offset of the union discriminant, in multiples of 16 bits.
	// 	discriminant = x.Uint16(capnp.DataOffset(x.DiscriminantOffset() * 2))
	// }
	// b/c of groups the field's index in this list is not necessarily exactly its ordinal.
	// however... if you want to identify a field by number, it may make the
	//most sense to use the field's index in this list rather than its ordinal.
	fields, e := x.Fields()
	if e != nil {
		t.Fatal(e)
	}
	for i := 0; i < fields.Len(); i++ {
		f := fields.At(i)
		n, _ := f.Name()
		t.Log(i, n)
		// CodeOrder
		// Annotations
		// DiscrimentValue
		// Fields are either "slots" or "groups"
		// - offset in units of the field size
		// - Type{} and Type_Which
		// - defaultValue
		// - hasExplicitDefault
	}
}

// https://github.com/capnproto/go-capnproto2/blob/master/std/capnp/schema.capnp
func hasDiscriminant(n schema.Node_structNode) bool {
	// Number of fields ( 0 or >1 ) in this struct which are members of a union.
	// If this is non-zero, then a 16-bit discriminant is present indicating the active member.
	return n.DiscriminantCount() > 0
}
