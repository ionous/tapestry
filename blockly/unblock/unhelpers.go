package unblock

import (
	"encoding/json"
	"errors"
	r "reflect"
	"strings"

	"git.sr.ht/~ionous/tapestry/lang/typeinfo"
	"git.sr.ht/~ionous/tapestry/support/inflect"
	"git.sr.ht/~ionous/tapestry/web/js"
)

type Creator interface {
	NewType(string) (any, bool)
}

func MakeBlockCreator(blocks []*typeinfo.TypeSet) Creator {
	m := make(typeMap)
	for _, b := range blocks {
		for _, ptr := range b.Signatures {
			t := ptr.TypeInfo()
			n := t.TypeName() // multiple signatures can generate the same type
			if _, ok := m[n]; !ok {
				m[n] = r.TypeOf(ptr).Elem()
			}
		}
	}
	return m
}

type typeMap map[string]r.Type

// returns .(typeinfo.Instance)
func (reg typeMap) NewType(name string) (ret any, okay bool) {
	if rtype, ok := reg[name]; ok {
		ret = r.New(rtype).Interface()
		okay = true
	}
	return
}

// fix: turn pascal to block case - TEST_NAME
func upper(n string) string {
	s := inflect.MixedCaseToSpaces(n)
	return strings.Replace(strings.ToUpper(s), " ", "_", -1)
}

func unstackName(n string) (ret string, okay bool) {
	const suffix = "_stack"
	if cnt := len(n); cnt > len(suffix) && n[0] == '_' && n[cnt-len(suffix):] == suffix {
		ret = n[1 : cnt-len(suffix)]
		okay = true
	}
	return
}

func readInput(m js.MapItem) (ret Input, err error) {
	if e := json.Unmarshal(m.Msg, &ret); e != nil {
		err = e
	} else if ret.BlockInfo == nil {
		err = errors.New("missing input")
	}
	return
}

func readValue(m js.MapItem) (ret any, err error) {
	err = json.Unmarshal(m.Msg, &ret)
	return
}
