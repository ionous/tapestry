package jsn

import (
	"bytes"
	"encoding/json"

	"github.com/ionous/errutil"
)

func Indent(str string) (ret string) {
	var indent bytes.Buffer
	if e := json.Indent(&indent, []byte(str), "", "  "); e != nil {
		ret = errutil.Sprint("indention error", e)
	} else {
		ret = indent.String()
	}
	return
}

func Compact(str string) (ret string) {
	var b bytes.Buffer
	if e := json.Compact(&b, []byte(str)); e != nil {
		ret = errutil.Sprint("compaction error", e)
	} else {
		ret = b.String()
	}
	return
}
