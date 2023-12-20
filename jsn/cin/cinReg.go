package cin

import (
	"hash/fnv"
	"io"
	r "reflect"

	"github.com/ionous/errutil"
)

// Generate instances of go types from story file signatures.
type TypeCreator interface {
	HasType(Hashed) bool
	NewFromSignature(Hashed) (any, error)
}

// hash to pointer to a nil value of the type
// ex. 15485098871275255450: (*Comparison)(nil)
type Signatures []map[uint64]any

// true if the signature is known
func (regs Signatures) HasType(sig Hashed) (okay bool) {
	_, okay = regs.FindType(sig)
	return
}

// return the nil pointer to the type.
func (regs Signatures) FindType(sig Hashed) (ret any, okay bool) {
	for _, reg := range regs {
		if a, ok := reg[sig.Value]; ok {
			ret, okay = a, ok
			break
		}
	}
	return
}

// create a new instance of the type and return its pointer
func (regs Signatures) NewFromSignature(sig Hashed) (ret any, err error) {
	if cmd, ok := regs.FindType(sig); !ok {
		err = errutil.New("unknown signature", sig)
	} else {
		ret = r.New(r.TypeOf(cmd).Elem()).Interface()
	}
	return
}

type Hashed struct {
	Value  uint64
	String string
}

// Hash helper for generating signatures lookups
// the slat includes all of the colon separated labels used to indicate a given type.
// the slot specifes the typename the slat fits into.
func Hash(slat, slot string) Hashed {
	str := slat
	if len(slot) > 0 {
		str = slot + "=" + slat
	}
	w := fnv.New64a()
	io.WriteString(w, str)
	return Hashed{Value: w.Sum64(), String: str}
}
