package cmdcompact

import (
	"git.sr.ht/~ionous/tapestry/support/files"
	"github.com/ionous/errutil"
)

type format int

const (
	jsonFormat format = iota
	tellFormat
)

func (f format) write(out string, data any, pretty bool) (err error) {
	