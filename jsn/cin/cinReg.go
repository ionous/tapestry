package cin

import (
	"hash/fnv"
	"io"
	r "reflect"

	"github.com/ionous/errutil"
)

type TypeCreator interface {
	HasType(string) bool
	NewFromSignature(string) (interface{}, error)
}

type Signatures []map[uint64]interface{}

func (regs Signatures) HasType(key string) (okay bool) {
	_, okay = regs.FindType(key)
	return
}

func (regs Signatures) FindType(key string) (ret interface{}, okay bool) {
	sig := Hash(key)
	for _, reg := range regs {
		if a, ok := reg[sig]; ok {
			ret, okay = a, ok
			break
		}
	}
	return
}

func (regs Signatures) NewFromSignature(key string) (ret interface{}, err error) {
	if cmd, ok := regs.FindType(key); !ok {
		err = errutil.New("unknown signature", key)
	} else {
		ret = r.New(r.TypeOf(cmd).Elem()).Interface()
	}
	return
}

// Hash helper for generating signatures lookups
func Hash(k string) uint64 {
	w := fnv.New64a()
	io.WriteString(w, k)
	return w.Sum64()
}
