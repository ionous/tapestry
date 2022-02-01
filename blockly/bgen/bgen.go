package bgen

import (
	"git.sr.ht/~ionous/tapestry/jsn"
	"git.sr.ht/~ionous/tapestry/jsn/chart"
	"git.sr.ht/~ionous/tapestry/web/js"
)

// convert the passed elements to blockly workspace format
func Convert(story jsn.Marshalee) (ret string, err error) {
	enc := chart.MakeEncoder()
	const header = `{"blocks": {"languageVersion": 0,"blocks": [`
	const footer = `]}}`
	var out js.Builder
	out.WriteString(header)
	if e := enc.Marshal(story, NewTopBlock(&enc, &out)); e != nil {
		err = e
	} else {
		out.WriteString(footer)
		ret = out.String()
	}
	return
}
