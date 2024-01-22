package block

import (
	"git.sr.ht/~ionous/tapestry/blockly/bconst"
	"git.sr.ht/~ionous/tapestry/lang/typeinfo"
	"git.sr.ht/~ionous/tapestry/web/js"
)

// convert the passed elements to blockly workspace format
func Convert(types bconst.Types, story typeinfo.Inspector) (ret string, err error) {
	const header = `{"blocks": {"languageVersion": 0,"blocks": [`
	const footer = `]}}`
	var out js.Builder
	out.WriteString(header)
	if e := Build(&out, story, types, true); e != nil {
		err = e
	} else {
		out.WriteString(footer)
		ret = out.String()
	}
	return
}
