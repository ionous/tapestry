package unblock

import (
	"encoding/json"
	"strings"

	inflect "git.sr.ht/~ionous/tapestry/inflect/en"
	"git.sr.ht/~ionous/tapestry/jsn"
	"git.sr.ht/~ionous/tapestry/web/js"
)

type TypeCreator interface {
	NewType(string) (any, bool)
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
		err = jsn.Missing
	}
	return
}

func readValue(m js.MapItem) (ret any, err error) {
	err = json.Unmarshal(m.Msg, &ret)
	return
}
