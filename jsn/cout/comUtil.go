package cout

import (
	"io"
	"strings"

	"git.sr.ht/~ionous/tapestry/jsn"
	"git.sr.ht/~ionous/tapestry/support/files"
)

func Marshal(cmd jsn.Marshalee, customFlow CustomFlow) (ret string, err error) {
	var out strings.Builder
	if e := marshal(&out, cmd, customFlow, false); e != nil {
		err = e
	} else {
		ret = out.String()
	}
	return
}

func MarshalIndent(out io.Writer, cmd jsn.Marshalee, customFlow CustomFlow) error {
	return marshal(out, cmd, customFlow, true)
}

func marshal(out io.Writer, cmd jsn.Marshalee, customFlow CustomFlow, indent bool) (err error) {
	if data, e := Encode(cmd, customFlow); e != nil {
		err = e
	} else {
		err = files.WriteJson(out, data, indent)
	}
	return
}
