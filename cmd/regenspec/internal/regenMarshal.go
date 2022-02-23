package regen

import (
	"encoding/json"
	"io"
	"strings"
)

func MarshalIndentOut(out io.Writer, i interface{}) {
	js := json.NewEncoder(out)
	js.SetEscapeHTML(false)
	js.SetIndent("", "  ")
	if e := js.Encode(i); e != nil {
		panic(e)
	}
}

func MarshalIndent(i interface{}) string {
	var b strings.Builder
	MarshalIndentOut(&b, i)
	return b.String()
}
