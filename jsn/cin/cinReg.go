package cin

import (
	"hash/fnv"
	"io"
	r "reflect"

	"github.com/ionous/errutil"
)

// Generate instances of go types from .if file signatures.
type TypeCreator interface {
	HasType(Hashed) bool
	NewFromSignature(Hashed) (interface{}, error)
}

type Signatures []map[uint64]interface{}

func (regs Signatures) HasType(sig Hashed) (okay bool) {
	_, okay = regs.FindType(sig)
	return
}

func (regs Signatures) FindType(sig Hashed) (ret interface{}, okay bool) {
	for _, reg := range regs {
		if a, ok := reg[sig.Value]; ok {
			ret, okay = a, ok
			break
		}
	}
	return
}

func (regs Signatures) NewFromSignature(sig Hashed) (ret interface{}, err error) {
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
